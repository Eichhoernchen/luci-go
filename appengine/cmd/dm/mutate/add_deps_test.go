// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package mutate

import (
	"testing"

	"github.com/luci/gae/impl/memory"
	"github.com/luci/gae/service/datastore"
	"github.com/luci/luci-go/appengine/cmd/dm/model"
	"github.com/luci/luci-go/common/api/dm/service/v1"
	. "github.com/luci/luci-go/common/testing/assertions"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestAddDeps(t *testing.T) {
	t.Parallel()

	Convey("AddDeps", t, func() {
		c := memory.Use(context.Background())

		ds := datastore.Get(c)

		aid := dm.NewAttemptID("quest", 1)
		a := model.MakeAttempt(c, aid)
		a.CurExecution = 1
		a.State = dm.Attempt_Executing
		ex := &model.Execution{
			ID: 1, Attempt: ds.KeyForObj(a), Token: []byte("sup"),
			State: dm.Execution_Running}

		ad := &AddDeps{
			Auth: &dm.Execution_Auth{
				Id:    dm.NewExecutionID("quest", 1, 1),
				Token: []byte("sup"),
			},
			ToAdd: dm.NewAttemptList(map[string][]uint32{
				"to":    {1, 2, 3},
				"top":   {1},
				"tp":    {1},
				"zebra": {17},
			}),
		}
		fds := model.FwdDepsFromList(c, aid, ad.ToAdd)

		Convey("Root", func() {
			So(ad.Root(c).String(), ShouldEqual, `dev~app::/Attempt,"quest|fffffffe"`)
		})

		Convey("RollForward", func() {

			Convey("Bad", func() {
				Convey("Bad ExecutionKey", func() {
					So(ds.PutMulti([]interface{}{a, ex}), ShouldBeNil)

					ad.Auth.Token = []byte("nerp")
					muts, err := ad.RollForward(c)
					So(err, ShouldBeRPCUnauthenticated, "execution Auth")
					So(muts, ShouldBeEmpty)
				})
			})

			Convey("Good", func() {
				So(ds.PutMulti([]interface{}{a, ex}), ShouldBeNil)

				Convey("All added already", func() {
					So(ds.PutMulti(fds), ShouldBeNil)

					muts, err := ad.RollForward(c)
					So(err, ShouldBeNil)
					So(muts, ShouldBeEmpty)
				})

				Convey("None added already", func() {
					muts, err := ad.RollForward(c)
					So(err, ShouldBeNil)
					So(len(muts), ShouldEqual, 2*len(fds))

					So(muts[0], ShouldResemble, &EnsureAttempt{&fds[0].Dependee})
					So(muts[1], ShouldResemble, &AddBackDep{
						Dep: fds[0].Edge(), NeedsAck: true})

					So(ds.Get(a), ShouldBeNil)
					So(ds.GetMulti(fds), ShouldBeNil)
					So(a.AddingDepsBitmap.Size(), ShouldEqual, len(fds))
					So(a.WaitingDepBitmap.Size(), ShouldEqual, len(fds))
					So(a.State, ShouldEqual, dm.Attempt_AddingDeps)
					So(fds[0].ForExecution, ShouldEqual, 1)
				})
			})
		})
	})
}
