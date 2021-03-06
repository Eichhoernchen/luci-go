// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package testconfig

import (
	"net/url"
	"testing"

	"github.com/luci/luci-go/common/config/impl/memory"
	configPB "github.com/luci/luci-go/common/proto/config"
	"github.com/luci/luci-go/luci_config/common/cfgtypes"
	"github.com/luci/luci-go/luci_config/server/cfgclient"
	"github.com/luci/luci-go/luci_config/server/cfgclient/backend"
	"github.com/luci/luci-go/luci_config/server/cfgclient/backend/client"
	"github.com/luci/luci-go/server/auth"
	"github.com/luci/luci-go/server/auth/authtest"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	. "github.com/smartystreets/goconvey/convey"
)

func tpb(msg proto.Message) string { return proto.MarshalTextString(msg) }

func accessCfg(access ...string) string {
	return tpb(&configPB.ProjectCfg{
		Access: access,
	})
}

func TestLocalService(t *testing.T) {
	t.Parallel()

	Convey(`Testing the local service`, t, func() {
		c := context.Background()

		fs := authtest.FakeState{
			Identity:       "user:foo@bar.baz",
			IdentityGroups: []string{"all", "special"},
		}
		c = auth.WithState(c, &fs)

		configs := map[string]memory.ConfigSet{
			"projects/foo": {
				"path.cfg":    "foo",
				"project.cfg": accessCfg("group:all"),
			},
			"projects/exclusive": {
				"path.cfg":    "exclusive",
				"project.cfg": accessCfg("group:special"),
			},
			"projects/exclusive/refs/heads/master": {
				"path.cfg": "exclusive master",
			},
			"projects/nouser": {
				"path.cfg":    "nouser",
				"project.cfg": accessCfg("group:impossible"),
			},
			"projects/nouser/refs/heads/master": {
				"path.cfg": "nouser master",
			},
			"services/baz": {
				"path.cfg": "service only",
			},
		}
		mbase := memory.New(configs)
		c = backend.WithBackend(c, &client.Backend{&Provider{
			Base: mbase,
		}})

		metaFor := func(configSet, path string) *cfgclient.Meta {
			cfg, err := mbase.GetConfig(c, configSet, path, false)
			if err != nil {
				panic(err)
			}
			return &cfgclient.Meta{
				ConfigSet:   cfgtypes.ConfigSet(cfg.ConfigSet),
				Path:        cfg.Path,
				ContentHash: cfg.ContentHash,
				Revision:    cfg.Revision,
			}
		}

		Convey(`Can get the service URL`, func() {
			So(cfgclient.ServiceURL(c), ShouldResemble, url.URL{Scheme: "test", Host: "example.com"})
		})

		Convey(`Can get a single config`, func() {
			var (
				val  string
				meta cfgclient.Meta
			)

			Convey(`AsService`, func() {
				So(cfgclient.Get(c, cfgclient.AsService, "projects/foo", "path.cfg", cfgclient.String(&val), &meta), ShouldBeNil)
				So(val, ShouldEqual, "foo")
				So(&meta, ShouldResemble, metaFor("projects/foo", "path.cfg"))

				So(cfgclient.Get(c, cfgclient.AsService, "projects/exclusive", "path.cfg", cfgclient.String(&val), nil), ShouldBeNil)
				So(val, ShouldEqual, "exclusive")

				So(cfgclient.Get(c, cfgclient.AsService, "services/baz", "path.cfg", cfgclient.String(&val), nil), ShouldBeNil)
				So(val, ShouldEqual, "service only")
			})

			Convey(`AsUser`, func() {
				So(cfgclient.Get(c, cfgclient.AsUser, "projects/foo", "path.cfg", cfgclient.String(&val), nil), ShouldBeNil)
				So(val, ShouldEqual, "foo")

				So(cfgclient.Get(c, cfgclient.AsUser, "projects/exclusive", "path.cfg", cfgclient.String(&val), nil), ShouldBeNil)
				So(val, ShouldEqual, "exclusive")

				So(cfgclient.Get(c, cfgclient.AsUser, "services/baz", "path.cfg", cfgclient.String(&val), nil),
					ShouldEqual, cfgclient.ErrNoConfig)
			})

			Convey(`AsAnonymous`, func() {
				fs.IdentityGroups = []string{"all"}

				So(cfgclient.Get(c, cfgclient.AsAnonymous, "projects/foo", "path.cfg", cfgclient.String(&val), nil), ShouldBeNil)
				So(val, ShouldEqual, "foo")

				So(cfgclient.Get(c, cfgclient.AsAnonymous, "projects/exclusive", "path.cfg", cfgclient.String(&val), nil),
					ShouldEqual, cfgclient.ErrNoConfig)
				So(cfgclient.Get(c, cfgclient.AsAnonymous, "services/baz", "path.cfg", cfgclient.String(&val), nil),
					ShouldEqual, cfgclient.ErrNoConfig)
			})
		})

		Convey(`Can get multiple configs`, func() {
			var vals []string
			var meta []*cfgclient.Meta

			Convey(`AsService`, func() {
				So(cfgclient.Projects(c, cfgclient.AsService, "path.cfg", cfgclient.StringSlice(&vals), &meta),
					ShouldBeNil)
				So(vals, ShouldResemble, []string{"exclusive", "foo", "nouser"})
				So(meta, ShouldResemble, []*cfgclient.Meta{
					metaFor("projects/exclusive", "path.cfg"),
					metaFor("projects/foo", "path.cfg"),
					metaFor("projects/nouser", "path.cfg"),
				})

				So(cfgclient.Refs(c, cfgclient.AsService, "path.cfg", cfgclient.StringSlice(&vals), &meta),
					ShouldBeNil)
				So(vals, ShouldResemble, []string{"exclusive master", "nouser master"})
				So(meta, ShouldResemble, []*cfgclient.Meta{
					metaFor("projects/exclusive/refs/heads/master", "path.cfg"),
					metaFor("projects/nouser/refs/heads/master", "path.cfg"),
				})
			})

			Convey(`AsUser`, func() {
				So(cfgclient.Projects(c, cfgclient.AsUser, "path.cfg", cfgclient.StringSlice(&vals), &meta),
					ShouldBeNil)
				So(vals, ShouldResemble, []string{"exclusive", "foo"})
				So(meta, ShouldResemble, []*cfgclient.Meta{
					metaFor("projects/exclusive", "path.cfg"),
					metaFor("projects/foo", "path.cfg"),
				})

				So(cfgclient.Refs(c, cfgclient.AsUser, "path.cfg", cfgclient.StringSlice(&vals), &meta),
					ShouldBeNil)
				So(vals, ShouldResemble, []string{"exclusive master"})
				So(meta, ShouldResemble, []*cfgclient.Meta{
					metaFor("projects/exclusive/refs/heads/master", "path.cfg"),
				})
			})

			Convey(`AsAnonymous`, func() {
				fs.IdentityGroups = []string{"all"}

				So(cfgclient.Projects(c, cfgclient.AsAnonymous, "path.cfg", cfgclient.StringSlice(&vals), &meta),
					ShouldBeNil)
				So(vals, ShouldResemble, []string{"foo"})
				So(meta, ShouldResemble, []*cfgclient.Meta{
					metaFor("projects/foo", "path.cfg"),
				})

				So(cfgclient.Refs(c, cfgclient.AsAnonymous, "path.cfg", cfgclient.StringSlice(&vals), &meta),
					ShouldBeNil)
				So(vals, ShouldHaveLength, 0)
				So(meta, ShouldHaveLength, 0)
			})
		})

		Convey(`Can get config set URLs`, func() {
			Convey(`AsService`, func() {
				u, err := cfgclient.GetConfigSetURL(c, cfgclient.AsService, "projects/foo")
				So(err, ShouldBeNil)
				So(u, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/fake-config/projects/foo"})

				u, err = cfgclient.GetConfigSetURL(c, cfgclient.AsService, "projects/exclusive")
				So(err, ShouldBeNil)
				So(u, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/fake-config/projects/exclusive"})

				u, err = cfgclient.GetConfigSetURL(c, cfgclient.AsService, "projects/nouser")
				So(err, ShouldBeNil)
				So(u, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/fake-config/projects/nouser"})
			})

			Convey(`AsUser`, func() {
				u, err := cfgclient.GetConfigSetURL(c, cfgclient.AsUser, "projects/foo")
				So(err, ShouldBeNil)
				So(u, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/fake-config/projects/foo"})

				u, err = cfgclient.GetConfigSetURL(c, cfgclient.AsUser, "projects/exclusive")
				So(err, ShouldBeNil)
				So(u, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/fake-config/projects/exclusive"})

				u, err = cfgclient.GetConfigSetURL(c, cfgclient.AsUser, "projects/nouser")
				So(err, ShouldEqual, cfgclient.ErrNoConfig)
			})

			Convey(`AsAnonymous`, func() {
				fs.IdentityGroups = []string{"all"}

				u, err := cfgclient.GetConfigSetURL(c, cfgclient.AsAnonymous, "projects/foo")
				So(err, ShouldBeNil)
				So(u, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/fake-config/projects/foo"})

				_, err = cfgclient.GetConfigSetURL(c, cfgclient.AsAnonymous, "projects/exclusive")
				So(err, ShouldEqual, cfgclient.ErrNoConfig)

				_, err = cfgclient.GetConfigSetURL(c, cfgclient.AsAnonymous, "projects/nouser")
				So(err, ShouldEqual, cfgclient.ErrNoConfig)
			})
		})
	})
}
