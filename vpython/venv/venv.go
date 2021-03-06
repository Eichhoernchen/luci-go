// Copyright 2017 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package venv

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/danjacques/gofslock/fslock"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"github.com/luci/luci-go/vpython/api/vpython"
	"github.com/luci/luci-go/vpython/python"
	"github.com/luci/luci-go/vpython/wheel"

	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/data/stringset"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/system/filesystem"
)

// EnvironmentVersion is an environment version string. It must advance each
// time the layout of a VirtualEnv environment changes.
const EnvironmentVersion = "v1"

// ErrNotComplete is a sentinel error returned by AssertCompleteAndLoad to
// indicate that the Environment is missing its completion flag.
var ErrNotComplete = errors.New("environment is not complete")

const (
	lockHeldDelay = 10 * time.Millisecond
)

// blocker is an fslock.Blocker implementation that sleeps lockHeldDelay in
// between attempts.
func blocker(c context.Context) fslock.Blocker {
	return func() error {
		logging.Debugf(c, "Lock is currently held. Sleeping %v and retrying...", lockHeldDelay)
		tr := clock.Sleep(c, lockHeldDelay)
		return tr.Err
	}
}

func withTempDir(l logging.Logger, prefix string, fn func(string) error) error {
	tdir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return errors.Annotate(err).Reason("failed to create temporary directory").Err()
	}
	defer func() {
		if err := filesystem.RemoveAll(tdir); err != nil {
			l.Infof("Failed to remove temporary directory: %s", err)
		}
	}()

	return fn(tdir)
}

// EnvRootFromStampPath calculates the environment root from an exported
// environment specification file path.
//
// The specification path is: <EnvRoot>/<SpecHash>/EnvironmentStampPath, so our
// EnvRoot is two directories up.
//
// We export EnvSpecPath as an asbolute path. However, since someone else
// could have overridden it or exported their own, let's make sure.
func EnvRootFromStampPath(path string) (string, error) {
	if err := filesystem.AbsPath(&path); err != nil {
		return "", errors.Annotate(err).
			Reason("failed to get absolute path for specification file path: %(path)s").
			Err()
	}
	return filepath.Dir(filepath.Dir(path)), nil
}

// With creates a new Env and executes "fn" with assumed ownership of that Env.
//
// The Context passed to "fn" will be cancelled if we lose perceived ownership
// of the configured environment. This is not an expected scenario, and should
// be considered an error condition. The Env passed to "fn" is valid only for
// the duration of the callback.
//
// It will lock around the VirtualEnv to ensure that multiple processes do not
// conflict with each other. If a VirtualEnv for this specification already
// exists, it will be used directly without any additional setup.
//
// If another process holds the lock, With will return an error (if blocking is
// false) or try again until it obtains the lock (if blocking is true).
func With(c context.Context, cfg Config, blocking bool, fn func(context.Context, *Env) error) error {
	// Track which VirtualEnv we use so we can exempt them from pruning.
	usedEnvs := stringset.New(2)

	// Start with an empty VirtualEnv. We will use this to probe the local
	// system.
	//
	// If our configured VirtualEnv is, itself, an empty then we can
	// skip this.
	var e *vpython.Environment
	if cfg.HasWheels() {
		// Use an empty VirtualEnv to probe the runtime environment.
		//
		// Disable pruning for this step since we'll be doing that later with the
		// full environment initialization.
		emptyEnv, err := cfg.WithoutWheels().makeEnv(c, nil)
		if err != nil {
			return errors.Annotate(err).Reason("failed to initialize empty probe environment").Err()
		}
		if err := emptyEnv.ensure(c, blocking); err != nil {
			return errors.Annotate(err).Reason("failed to create empty probe environment").Err()
		}

		usedEnvs.Add(emptyEnv.Name)
		e = emptyEnv.Environment
	}

	// Run the real config, now with runtime data.
	env, err := cfg.makeEnv(c, e)
	if err != nil {
		return err
	}
	usedEnvs.Add(env.Name)
	return env.withImpl(c, blocking, usedEnvs, fn)
}

// Delete removes all resources consumed by an environment.
//
// Delete will acquire an exclusive lock on the environment and assert that it
// is not in use prior to deletion. This is non-blocking, and an error will be
// returned if the lock could not be acquired.
//
// If deletion fails, a wrapped error will be returned.
func Delete(c context.Context, cfg Config) error {
	// Attempt to acquire the environment's lock.
	e, err := cfg.makeEnv(c, nil)
	if err != nil {
		return err
	}
	return e.Delete(c)
}

// Env is a fully set-up Python virtual environment. It is configured
// based on the contents of an vpython.Spec file by Setup.
//
// Env should not be instantiated directly; it must be created by calling
// Config.Env.
//
// All paths in Env are absolute.
type Env struct {
	// Config is this Env's Config, fully-resolved.
	Config *Config

	// Root is the Env container's root directory path.
	Root string

	// Name is the hash of the specification file for this Env.
	Name string

	// Python is the path to the Env Python interpreter.
	Python string

	// Environment is the resolved Python environment that this VirtualEnv is
	// configured with. It is resolved during environment setup and saved into the
	// environment's EnvironmentStampPath.
	Environment *vpython.Environment

	// BinDir is the VirtualEnv "bin" directory, containing Python and installed
	// wheel binaries.
	BinDir string
	// EnvironmentStampPath is the path to the vpython.Environment stamp file that
	// details a constructed environment. It will be in text protobuf format, and,
	// therefore, suitable for input to other "vpython" invocations.
	EnvironmentStampPath string

	// lockPath is the path to this Env-specific lock file. It will be at:
	// "<baseDir>/.<name>.lock".
	lockPath string
	// completeFlagPath is the path to this Env's complete flag.
	// It will be at "<Root>/complete.flag".
	completeFlagPath string
	// interpreter is the VirtualEnv Python interpreter.
	//
	// It is configured to use the Python member as its base executable, and is
	// initialized on first call to Interpreter.
	interpreter *python.Interpreter
}

// ensure ensures that the configured VirtualEnv is set-up. This may involve
// a fast-path completion flag check or a slow lock/build phase.
func (e *Env) ensure(c context.Context, blocking bool) (err error) {
	// Fastest path: If this environment is already complete, then there is no
	// additional setup necessary.
	if err := e.AssertCompleteAndLoad(); err == nil {
		logging.Debugf(c, "Environment is already initialized: %s", e.Environment)
		return nil
	}

	// Repeatedly try and create our Env. We do this so that if we
	// encounter a lock, we will let the other process finish and try and leverage
	// its success.
	for {
		// We will be creating the Env. Acquire an exclusive lock so that any other
		// processes will wait for our setup to complete.
		switch lock, err := e.acquireExclusiveLock(); err {
		case nil:
			// We MUST successfully release our exclusive lock on completion.
			return mustReleaseLock(c, lock, func() error {
				// Fast path: if our complete flag is present, assume that the
				// environment is setup and complete. No additional work is necessary.
				err := e.AssertCompleteAndLoad()
				if err == nil {
					// This will generally happen if another process created the environment
					// in between when we checked for the completion stamp initially and
					// when we actually obtained the lock.
					logging.Debugf(c, "Completion flag found! Environment is set-up: %s", e.completeFlagPath)
					return nil
				}
				logging.WithError(err).Debugf(c, "VirtualEnv is not complete.")

				// No complete flag. Create a new VirtualEnv here.
				if err := e.createLocked(c); err != nil {
					return errors.Annotate(err).Reason("failed to create new VirtualEnv").Err()
				}

				// Mark that this environment is complete. This MUST succeed so other
				// instances know that this environment is complete.
				if err := e.touchCompleteFlagLocked(); err != nil {
					return errors.Annotate(err).Reason("failed to create complete flag").Err()
				}

				logging.Debugf(c, "Successfully created new virtual environment [%s]!", e.Name)
				return nil
			})

		case fslock.ErrLockHeld:
			// We couldn't get an exclusive lock. Try again to load the environment
			// stamp, asserting the existence of the completion flag in the process.
			// If another process has created the environment, we may be able to use
			// it without ever having to obtain its lock!
			if err := e.AssertCompleteAndLoad(); err == nil {
				logging.Infof(c, "Environment was completed while waiting for lock: %s", e.EnvironmentStampPath)
				return nil
			}

			logging.Fields{
				logging.ErrorKey: err,
				"path":           e.EnvironmentStampPath,
			}.Debugf(c, "Lock is held, and environment is not complete.")
			if !blocking {
				return errors.Annotate(err).Reason("VirtualEnv lock is currently held (non-blocking)").Err()
			}

			// Some other process holds the lock. Sleep a little and retry.
			if err := blocker(c)(); err != nil {
				return err
			}

		default:
			return errors.Annotate(err).Reason("failed to create VirtualEnv").Err()
		}
	}
}

func (e *Env) withImpl(c context.Context, blocking bool, used stringset.Set,
	fn func(context.Context, *Env) error) (err error) {

	// Setup the VirtualEnv environment.
	//
	// Setup will obtain an exclusive lock on the environment for set-up and
	// release it when it finishes. We will then obtain a shared lock on the
	// environment to represent its use.
	//
	// An exclusive lock may be taken on the environment in between the exclusive
	// lock being released and the shared lock being obtained. If that happens,
	// we will (if configured to block) repeat our loop and call "ensure" again.
	for {
		if err := e.ensure(c, blocking); err != nil {
			return err
		}

		// Acquire a shared lock on the environment to note its continued usage.
		switch lock, err := e.acquireSharedLock(); err {
		case nil:
			logging.Debugf(c, "Acquired shared lock for: %s", e.Name)

			environmentWasIncomplete := false
			err := mustReleaseLock(c, lock, func() error {
				// Assert that the environment wasn't deleted in between creation and
				// our acquisition of the lock.
				if err := e.assertComplete(); err != nil {
					logging.WithError(err).Infof(c, "Environment is no longer complete; recreating...")
					environmentWasIncomplete = true
					return nil
				}

				// Try and touch the complete flag to update its timestamp and mark this
				// environment's utility.
				if err := e.touchCompleteFlagLocked(); err != nil {
					logging.Debugf(c, "Failed to update environment timestamp.")
				}

				// Perform a pruning round. Failure is non-fatal.
				if perr := prune(c, e.Config, used); perr != nil {
					logging.WithError(perr).Infof(c, "Failed to perform pruning round after initialization.")
				}

				return fn(c, e)
			})
			if err != nil {
				return err
			}

			logging.Debugf(c, "Released shared lock for: %s", e.Name)
			if !environmentWasIncomplete {
				return nil
			}

		case fslock.ErrLockHeld:
			logging.Fields{
				logging.ErrorKey: err,
				"path":           e.EnvironmentStampPath,
			}.Debugf(c, "Could not obtain shared usage lock.")
			if !blocking {
				return errors.Annotate(err).Reason("VirtualEnv lock is currently held (non-blocking)").Err()
			}

			// Some other process holds the lock. Sleep a little and retry.
			if err := blocker(c)(); err != nil {
				return err
			}

		default:
			return errors.Annotate(err).Reason("failed to use VirtualEnv").Err()
		}
	}
}

// Interpreter returns the VirtualEnv's isolated Python Interpreter instance.
func (e *Env) Interpreter() *python.Interpreter {
	if e.interpreter == nil {
		e.interpreter = &python.Interpreter{
			Python: e.Python,
		}
	}
	return e.interpreter
}

func (e *Env) acquireExclusiveLock() (fslock.Handle, error) { return fslock.Lock(e.lockPath) }

func (e *Env) acquireSharedLock() (fslock.Handle, error) { return fslock.LockShared(e.lockPath) }

func (e *Env) withExclusiveLockNonBlocking(fn func() error) error {
	return fslock.With(e.lockPath, fn)
}

// WriteEnvironmentStamp writes a text protobuf form of spec to path.
func (e *Env) WriteEnvironmentStamp() error {
	environment := e.Environment
	if environment == nil {
		environment = &vpython.Environment{}
	}
	return writeTextProto(e.EnvironmentStampPath, environment)
}

// AssertCompleteAndLoad asserts that the VirtualEnv's completion
// flag exists. If it does, the environment's stamp is loaded into e.Environment
// and nil is returned.
//
// An error is returned if the completion flag does not exist, or if the
// VirtualEnv environment stamp could not be loaded.
func (e *Env) AssertCompleteAndLoad() error {
	if err := e.assertComplete(); err != nil {
		return err
	}

	content, err := ioutil.ReadFile(e.EnvironmentStampPath)
	if err != nil {
		return errors.Annotate(err).Reason("failed to load file from: %(path)s").
			D("path", e.EnvironmentStampPath).
			Err()
	}

	var environment vpython.Environment
	if err := proto.UnmarshalText(string(content), &environment); err != nil {
		return errors.Annotate(err).Reason("failed to unmarshal vpython.Env stamp from: %(path)s").
			D("path", e.EnvironmentStampPath).
			Err()
	}
	e.Environment = &environment
	return nil
}

func (e *Env) assertComplete() error {
	// Ensure that the environment has its completion flag.
	switch _, err := os.Stat(e.completeFlagPath); {
	case filesystem.IsNotExist(err):
		return ErrNotComplete
	case err != nil:
		return errors.Annotate(err).Reason("failed to check for completion flag").Err()
	default:
		return nil
	}
}

func (e *Env) createLocked(c context.Context) error {
	// If our root directory already exists, delete it.
	if _, err := os.Stat(e.Root); err == nil {
		logging.Infof(c, "Deleting existing VirtualEnv: %s", e.Root)
		if err := filesystem.RemoveAll(e.Root); err != nil {
			return errors.Reason("failed to remove existing root").Err()
		}
	}

	// Make sure our environment's base directory exists.
	if err := filesystem.MakeDirs(e.Root); err != nil {
		return errors.Annotate(err).Reason("failed to create environment root").Err()
	}
	logging.Infof(c, "Using virtual environment root: %s", e.Root)

	// Build our package list. Always install our base VirtualEnv package.
	packages := make([]*vpython.Spec_Package, 1, 1+len(e.Environment.Spec.Wheel))
	packages[0] = e.Environment.Spec.Virtualenv
	packages = append(packages, e.Environment.Spec.Wheel...)

	// Create a directory to bootstrap VirtualEnv from.
	//
	// This directory will be a very short-named temporary directory. This is
	// because it really quickly runs up into traditional Windows path limitations
	// when ZIP-importing sub-sub-sub-sub-packages (e.g., pip, requests, etc.).
	//
	// We will clean this directory up on termination.
	err := withTempDir(logging.Get(c), "vpython_bootstrap", func(bootstrapDir string) error {
		pkgDir := filepath.Join(bootstrapDir, "packages")
		if err := filesystem.MakeDirs(pkgDir); err != nil {
			return errors.Annotate(err).Reason("could not create bootstrap packages directory").Err()
		}

		if err := e.downloadPackages(c, pkgDir, packages); err != nil {
			return errors.Annotate(err).Reason("failed to download packages").Err()
		}

		// Installing base VirtualEnv.
		if err := e.installVirtualEnv(c, pkgDir); err != nil {
			return errors.Annotate(err).Reason("failed to install VirtualEnv").Err()
		}

		// Load PEP425 tags, if we don't already have them.
		if e.Environment.Pep425Tag == nil {
			pep425Tags, err := e.getPEP425Tags(c)
			if err != nil {
				return errors.Annotate(err).Reason("failed to get PEP425 tags").Err()
			}
			e.Environment.Pep425Tag = pep425Tags
		}

		// Install our wheel files.
		if len(e.Environment.Spec.Wheel) > 0 {
			// Install wheels into our VirtualEnv.
			if err := e.installWheels(c, bootstrapDir, pkgDir); err != nil {
				return errors.Annotate(err).Reason("failed to install wheels").Err()
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Write our specification file.
	if err := e.WriteEnvironmentStamp(); err != nil {
		return errors.Annotate(err).Reason("failed to write environment stamp file to: %(path)s").
			D("path", e.EnvironmentStampPath).
			Err()
	}
	logging.Debugf(c, "Wrote environment stamp file to: %s", e.EnvironmentStampPath)

	// Finalize our VirtualEnv for bootstrap execution.
	if err := e.finalize(c); err != nil {
		return errors.Annotate(err).Reason("failed to prepare VirtualEnv").Err()
	}

	return nil
}

func (e *Env) downloadPackages(c context.Context, dst string, packages []*vpython.Spec_Package) error {
	// Create a wheel sub-directory underneath of root.
	logging.Debugf(c, "Loading %d package(s) into: %s", len(packages), dst)
	if err := e.Config.Loader.Ensure(c, dst, packages); err != nil {
		return errors.Annotate(err).Reason("failed to download packages").Err()
	}
	return nil
}

func (e *Env) installVirtualEnv(c context.Context, pkgDir string) error {
	// Create our VirtualEnv package staging sub-directory underneath of root.
	bsDir := filepath.Join(e.Root, ".virtualenv")
	if err := filesystem.MakeDirs(bsDir); err != nil {
		return errors.Annotate(err).Reason("failed to create VirtualEnv bootstrap directory").
			D("path", bsDir).
			Err()
	}

	// Identify the virtualenv directory: will have "virtualenv-" prefix.
	matches, err := filepath.Glob(filepath.Join(pkgDir, "virtualenv-*"))
	if err != nil {
		return errors.Annotate(err).Reason("failed to glob for 'virtualenv-' directory").Err()
	}
	if len(matches) == 0 {
		return errors.Reason("no 'virtualenv-' directory provided by package").Err()
	}

	logging.Debugf(c, "Creating VirtualEnv at: %s", e.Root)
	cmd := e.Config.systemInterpreter().IsolatedCommand(c,
		"virtualenv.py",
		"--no-download",
		e.Root)
	cmd.Dir = matches[0]
	attachOutputForLogging(c, logging.Debug, cmd)
	if err := cmd.Run(); err != nil {
		return errors.Annotate(err).Reason("failed to create VirtualEnv").Err()
	}

	return nil
}

// getPEP425Tags calls Python's pip.pep425tags package to retrieve the tags.
//
// This must be run while "pip" is installed in the VirtualEnv.
func (e *Env) getPEP425Tags(c context.Context) ([]*vpython.Pep425Tag, error) {
	// This script will return a list of 3-entry lists:
	// [0]: version (e.g., "cp27")
	// [1]: abi (e.g., "cp27mu", "none")
	// [2]: arch (e.g., "x86_64", "armv7l", "any")
	const script = `import json;` +
		`import pip.pep425tags;` +
		`import sys;` +
		`sys.stdout.write(json.dumps(pip.pep425tags.get_supported()))`
	type pep425TagEntry []string

	cmd := e.Interpreter().IsolatedCommand(c, "-c", script)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	attachOutputForLogging(c, logging.Debug, cmd)
	if err := cmd.Run(); err != nil {
		return nil, errors.Annotate(err).Reason("failed to get PEP425 tags").Err()
	}

	var tagEntries []pep425TagEntry
	if err := json.Unmarshal(stdout.Bytes(), &tagEntries); err != nil {
		return nil, errors.Annotate(err).Reason("failed to unmarshal PEP425 tag output: %(output)s").
			D("output", stdout.String()).
			Err()
	}

	tags := make([]*vpython.Pep425Tag, len(tagEntries))
	for i, te := range tagEntries {
		if len(te) != 3 {
			return nil, errors.Reason("invalid PEP425 tag entry: %(entry)v").
				D("entry", te).
				D("index", i).
				Err()
		}

		tags[i] = &vpython.Pep425Tag{
			Version: te[0],
			Abi:     te[1],
			Arch:    te[2],
		}
	}

	// If we're Debug-logging, calculate and display the tags that were probed.
	if logging.IsLogging(c, logging.Debug) {
		tagStr := make([]string, len(tags))
		for i, t := range tags {
			tagStr[i] = t.TagString()
		}
		logging.Debugf(c, "Loaded PEP425 tags: [%s]", strings.Join(tagStr, ", "))
	}

	return tags, nil
}

func (e *Env) installWheels(c context.Context, bootstrapDir, pkgDir string) error {
	// Identify all downloaded wheels and parse them.
	wheels, err := wheel.ScanDir(pkgDir)
	if err != nil {
		return errors.Annotate(err).Reason("failed to load wheels").Err()
	}

	// Build a "wheel" requirements file.
	reqPath := filepath.Join(bootstrapDir, "requirements.txt")
	logging.Debugf(c, "Rendering requirements file to: %s", reqPath)
	if err := wheel.WriteRequirementsFile(reqPath, wheels); err != nil {
		return errors.Annotate(err).Reason("failed to render requirements file").Err()
	}

	cmd := e.Interpreter().IsolatedCommand(c,
		"-m", "pip",
		"install",
		"--use-wheel",
		"--compile",
		"--no-index",
		"--find-links", pkgDir,
		"--requirement", reqPath)
	attachOutputForLogging(c, logging.Debug, cmd)
	if err := cmd.Run(); err != nil {
		return errors.Annotate(err).Reason("failed to install wheels").Err()
	}
	return nil
}

func (e *Env) finalize(c context.Context) error {
	// Change all files to read-only, except:
	// - Our root directory, which must be writable in order to update our
	//   completion flag.
	// - Our environment stamp, which must be trivially re-writable.
	if !e.Config.testLeaveReadWrite {
		err := filesystem.MakeReadOnly(e.Root, func(path string) bool {
			switch path {
			case e.Root, e.completeFlagPath:
				return false
			default:
				return true
			}
		})
		if err != nil {
			return errors.Annotate(err).Reason("failed to mark environment read-only").Err()
		}
	}
	return nil
}

func (e *Env) touchCompleteFlagLocked() error {
	if err := filesystem.Touch(e.completeFlagPath, time.Time{}, 0644); err != nil {
		return errors.Annotate(err).Err()
	}
	return nil
}

// Delete removes all resources consumed by an environment.
//
// Delete will acquire an exclusive lock on the environment and assert that it
// is not in use prior to deletion. This is non-blocking, and an error will be
// returned if the lock could not be acquired.
//
// If the environment was not deleted, a non-nil wrapped error will be returned.
// If the deletion failed because the lock was held, a wrapped
// fslock.ErrLockHeld  will be returned.
func (e *Env) Delete(c context.Context) error {
	removedLock := false
	err := e.withExclusiveLockNonBlocking(func() error {
		logging.Debugf(c, "(Delete) Got exclusive lock for: %s", e.Name)

		// Delete our environment directory.
		if err := filesystem.RemoveAll(e.Root); err != nil {
			return errors.Annotate(err).Reason("failed to delete environment root").Err()
		}

		// Attempt to delete our lock. On POSIX systems, this will successfully
		// delete the lock. On Windows systems, there will be contention, since the
		// lock is held by us.
		//
		// In this case, we'll try again to delete it after we release the lock.
		// If someone else takes out the lock in between, the delete will similarly
		// fail.
		if err := os.Remove(e.lockPath); err == nil {
			removedLock = true
		} else {
			logging.WithError(err).Debugf(c, "failed to delete lock file while holding lock: %s", e.lockPath)
		}
		return nil
	})
	logging.Debugf(c, "(Delete) Released exclusive lock for: %s", e.Name)

	if err != nil {
		return errors.Annotate(err).Reason("failed to delete environment").Err()
	}

	// Try and remove the lock now that we don't hold it.
	if !removedLock {
		if err := os.Remove(e.lockPath); err != nil {
			return errors.Annotate(err).Reason("failed to remove lock").Err()
		}
	}
	return nil
}

// completionFlagTimestamp returns the timestamp on the environment's completion
// flag.
//
// If the completion flag does not exist (incomplete environment), a zero time
// will be returned with no error.
//
// If an error is encountered while checking for the timestamp, it will be
// returned.
func (e *Env) completionFlagTimestamp() (time.Time, error) {
	// Read the complete flag file's timestamp.
	switch st, err := os.Stat(e.completeFlagPath); {
	case err == nil:
		return st.ModTime(), nil

	case os.IsNotExist(err):
		return time.Time{}, nil

	default:
		return time.Time{}, errors.Annotate(err).Reason("failed to stat completion flag: %(path)s").
			D("path", e.completeFlagPath).
			Err()
	}
}

func attachOutputForLogging(c context.Context, l logging.Level, cmd *exec.Cmd) {
	if logging.IsLogging(c, logging.Info) {
		logging.Infof(c, "Running Python command (cwd=%s): %s",
			cmd.Dir, strings.Join(cmd.Args, " "))
	}

	if logging.IsLogging(c, l) {
		if cmd.Stdout == nil {
			cmd.Stdout = os.Stdout
		}
		if cmd.Stderr == nil {
			cmd.Stderr = os.Stderr
		}
	}
}

// mustReleaseLock calls the wrapped function, releasing the lock at the end
// of its execution. If the lock could not be released, this function will
// panic, since the locking state can no longer be determined.
func mustReleaseLock(c context.Context, lock fslock.Handle, fn func() error) error {
	defer func() {
		if err := lock.Unlock(); err != nil {
			errors.Log(c, errors.Annotate(err).Reason("failed to release lock").Err())
			panic(err)
		}
	}()
	return fn()
}
