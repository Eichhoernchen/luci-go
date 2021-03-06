// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/errors"
	log "github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/luci_config/common/cfgtypes"
	"github.com/luci/luci-go/luci_config/server/cfgclient"

	"golang.org/x/net/context"
)

// ChangePoller polls a configuration files for changes. If it changes,
// the OnChange function will be called and the polling will stop.
type ChangePoller struct {
	// ConfigSet is the slice of config paths to watch.
	ConfigSet cfgtypes.ConfigSet
	// Path is the path of the config to watch.
	Path string

	// Period is the amount of time in between checks to see if the configuration
	// has been updated. If <= 0, the poller will refrain from polling, and Run
	// will immediately exit.
	Period time.Duration
	// OnChange is the function that will be called if a configuration change has
	// been observed.
	//
	// Polling will be blocked until OnChange returns. If the Context supplied to
	// Run is canceled by OnChange, Run will exit at the beginning of the next
	// poll round.
	OnChange func()

	// ContentHash is the config's hash. This should be set to the initial config
	// value, either directly or via a Refresh call, before Run is called.
	ContentHash string
}

// Run starts polling for changes. It will stop when the Context is cancelled.
func (p *ChangePoller) Run(c context.Context) {
	if p.Period <= 0 {
		return
	}

	for {
		// If our Context has been canceled, terminate.
		select {
		case <-c.Done():
			log.WithError(c.Err()).Warningf(c, "Terminating poll loop: context has been cancelled.")
			return
		default:
			// Continue
		}

		log.Fields{
			"timeout": p.Period,
		}.Debugf(c, "Entering change check poll loop...")
		if tr := clock.Sleep(c, p.Period); tr.Incomplete() {
			log.WithError(tr.Err).Debugf(c, "Context cancelled, shutting down change poller.")
			return
		}

		log.Infof(c, "Change check timeout triggered, checking configuration...")
		lastHash := p.ContentHash
		switch err := p.Refresh(c); {
		case err != nil:
			log.WithError(err).Errorf(c, "Failed to refresh config.")

		case lastHash != p.ContentHash:
			log.Fields{
				"originalHash": lastHash,
				"newHash":      p.ContentHash,
			}.Warningf(c, "Configuration content hash has changed.")
			if p.OnChange != nil {
				p.OnChange()
			}

		default:
			log.Fields{
				"currentHash": lastHash,
			}.Debugf(c, "Content hash matches.")
		}
	}
}

// Refresh reloads the configuration value, updating ContentHash.
func (p *ChangePoller) Refresh(c context.Context) error {
	var meta cfgclient.Meta
	if err := cfgclient.Get(c, cfgclient.AsService, p.ConfigSet, p.Path, nil, &meta); err != nil {
		return errors.Annotate(err).Reason("failed to reload config %(configSet)s :: %(path)s").
			D("configSet", p.ConfigSet).
			D("path", p.Path).
			Err()
	}

	p.ContentHash = meta.ContentHash
	return nil
}
