// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package deps

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/luci/gae/service/datastore"
	"github.com/luci/luci-go/appengine/cmd/dm/model"
	"github.com/luci/luci-go/appengine/tumble"
	"github.com/luci/luci-go/common/api/dm/service/v1"
	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/clock/testclock"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/logging/memlogger"
	google_pb "github.com/luci/luci-go/common/proto/google"
	. "github.com/luci/luci-go/common/testing/assertions"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func ensureQuest(c context.Context, name string) string {
	qsts, err := (&deps{}).EnsureQuests(c, &dm.EnsureQuestsReq{
		ToEnsure: []*dm.Quest_Desc{{
			DistributorConfigName: "foof",
			JsonPayload:           fmt.Sprintf(`{"name": "%s"}`, name),
		}},
	})
	So(err, ShouldBeNil)
	return qsts.QuestIds[0].Id
}

func execute(c context.Context, aid *dm.Attempt_ID) {
	// takes an NeedsExecution attempt, and moves it to Executing
	err := datastore.Get(c).RunInTransaction(func(c context.Context) error {
		ds := datastore.Get(c)
		atmpt := &model.Attempt{ID: *aid}
		if err := ds.Get(atmpt); err != nil {
			panic(err)
		}

		atmpt.CurExecution++
		atmpt.MustModifyState(c, dm.Attempt_Executing)

		So(ds.PutMulti([]interface{}{atmpt, &model.Execution{
			ID:      atmpt.CurExecution,
			Created: clock.Now(c),
			State:   dm.Execution_Running,
			Attempt: ds.KeyForObj(atmpt),
			Token:   []byte("sekret"),
		}}), ShouldBeNil)
		return nil
	}, nil)
	So(err, ShouldBeNil)
}

func depOn(c context.Context, from *dm.Attempt_ID, to ...*dm.Attempt_ID) {
	req := &dm.AddDepsReq{
		Auth: &dm.Execution_Auth{
			Id:    dm.NewExecutionID(from.Quest, from.Id, 1),
			Token: []byte("sekret")},
		Deps: dm.NewAttemptFanout(nil),
	}
	req.Deps.AddAIDs(to...)

	rsp, err := (&deps{}).AddDeps(c, req)
	So(err, ShouldBeNil)
	So(rsp.ShouldHalt, ShouldBeTrue)
}

func walkShouldReturn(c context.Context, log bool) func(request interface{}, expect ...interface{}) string {
	normalize := func(gd *dm.GraphData) *dm.GraphData {
		data, err := proto.Marshal(gd)
		if err != nil {
			panic(err)
		}
		So(err, ShouldBeNil)
		ret := &dm.GraphData{}
		if err := proto.Unmarshal(data, ret); err != nil {
			panic(err)
		}
		return ret
	}

	return func(request interface{}, expect ...interface{}) string {
		r := request.(*dm.WalkGraphReq)
		e := expect[0].(*dm.GraphData)
		if log {
			logging.Get(c).(*memlogger.MemLogger).Reset()
		}
		ret, err := (&deps{}).WalkGraph(c, r)
		if nilExpect := ShouldErrLike(err, nil); nilExpect != "" {
			return nilExpect
		}
		if log {
			for _, msg := range logging.Get(c).(*memlogger.MemLogger).Messages() {
				Println(msg)
			}
		}
		return ShouldResemble(normalize(ret), e)
	}
}

func WalkShouldReturn(c context.Context) func(request interface{}, expect ...interface{}) string {
	return walkShouldReturn(c, false)
}

func WalkShouldReturnDbg(c context.Context) func(request interface{}, expect ...interface{}) string {
	return walkShouldReturn(c, true)
}

func TestWalkGraph(t *testing.T) {
	t.Parallel()

	Convey("WalkGraph", t, func() {
		ttest := tumble.Testing{}
		c := ttest.Context()

		ds := datastore.Get(c)
		s := &deps{}

		req := &dm.WalkGraphReq{
			Queries: []*dm.GraphQuery{
				dm.AttemptListQueryL(map[string][]uint32{"quest": {1}})},
			Limit: &dm.WalkGraphReq_Limit{MaxDepth: 1},
		}
		So(req.Normalize(), ShouldBeNil)

		Convey("no attempt", func() {
			So(req, WalkShouldReturn(c), &dm.GraphData{
				Quests: map[string]*dm.Quest{"quest": {
					Attempts: map[uint32]*dm.Attempt{
						1: {DNE: true},
					},
				}},
			})
		})

		Convey("good", func() {
			ds.Testable().Consistent(true)

			wDesc := &dm.Quest_Desc{
				DistributorConfigName: "foof",
				JsonPayload:           `{"name":"w"}`,
			}
			w := ensureQuest(c, "w")
			//qid := dm.NewQuestID(w)
			aid := dm.NewAttemptID(w, 1)
			eid := dm.NewExecutionID(w, 1, 1)

			_, err := s.EnsureAttempt(c, &dm.EnsureAttemptReq{ToEnsure: aid})
			So(err, ShouldBeNil)

			req.Queries[0].GetAttemptList().Attempt = dm.NewAttemptFanout(
				map[string][]uint32{w: {1}})

			Convey("include nothing", func() {
				So(req, WalkShouldReturn(c), &dm.GraphData{
					Quests: map[string]*dm.Quest{
						w: {
							Attempts: map[uint32]*dm.Attempt{1: {}},
						},
					},
				})
			})

			Convey("quest dne", func() {
				req.Include.QuestData = true
				req.Limit.MaxDepth = 1
				req.Queries[0].GetAttemptList().Attempt = dm.NewAttemptFanout(
					map[string][]uint32{"noex": {1}})
				So(req, WalkShouldReturn(c), &dm.GraphData{
					Quests: map[string]*dm.Quest{
						"noex": {
							DNE:      true,
							Attempts: map[uint32]*dm.Attempt{1: {DNE: true}},
						},
					},
				})
			})

			Convey("no dependencies", func() {
				req.Include.AttemptData = true
				req.Include.QuestData = true
				req.Include.NumExecutions = 128

				now := datastore.RoundTime(clock.Now(c))
				aExpect := dm.NewAttemptNeedsExecution(now)
				aExpect.Data.Created = google_pb.NewTimestamp(now)
				aExpect.Data.Modified = google_pb.NewTimestamp(now)
				So(req, WalkShouldReturn(c), &dm.GraphData{
					Quests: map[string]*dm.Quest{
						w: {
							Data: &dm.Quest_Data{
								Created: google_pb.NewTimestamp(now),
								Desc:    wDesc,
							},
							Attempts: map[uint32]*dm.Attempt{1: aExpect},
						},
					},
				})
			})

			Convey("finished", func() {
				execute(c, aid)

				_, err := s.FinishAttempt(c, &dm.FinishAttemptReq{
					Auth: &dm.Execution_Auth{
						Id:    eid,
						Token: []byte("sekret"),
					},
					JsonResult: `{"data": ["very", "yes"]}`,
					Expiration: google_pb.NewTimestamp(clock.Now(c).Add(time.Hour * 24 * 4)),
				})
				So(err, ShouldBeNil)

				ex := &model.Execution{ID: 1, Attempt: ds.MakeKey("Attempt", aid.DMEncoded())}
				So(ds.Get(ex), ShouldBeNil)
				ex.State = dm.Execution_Finished
				So(ds.Put(ex), ShouldBeNil)

				req.Include.AttemptData = true
				req.Include.AttemptResult = true
				req.Include.NumExecutions = 128
				aExpect := dm.NewAttemptFinished(
					datastore.RoundTime(clock.Now(c).Add(time.Hour*24*4)),
					`{"data":["very","yes"]}`)
				aExpect.Data.Created = google_pb.NewTimestamp(datastore.RoundTime(clock.Now(c)))
				aExpect.Data.Modified = google_pb.NewTimestamp(datastore.RoundTime(clock.Now(c)))
				aExpect.Data.NumExecutions = 1
				aExpect.Executions = map[uint32]*dm.Execution{1: {
					State: dm.Execution_Finished,
					Data:  &dm.Execution_Data{Created: google_pb.NewTimestamp(datastore.RoundTime(clock.Now(c)))},
				}}
				So(req, WalkShouldReturn(c), &dm.GraphData{
					Quests: map[string]*dm.Quest{
						w: {
							Attempts: map[uint32]*dm.Attempt{1: aExpect},
						},
					},
				})
			})

			Convey("attemptRange", func() {
				x := ensureQuest(c, "x")
				execute(c, dm.NewAttemptID(w, 1))
				depOn(c, dm.NewAttemptID(w, 1), dm.NewAttemptID(x, 1), dm.NewAttemptID(x, 2), dm.NewAttemptID(x, 3), dm.NewAttemptID(x, 4))
				ttest.Drain(c)

				req.Limit.MaxDepth = 1
				Convey("normal", func() {
					req.Queries[0] = dm.AttemptRangeQuery(x, 2, 4)
					So(req, WalkShouldReturn(c), &dm.GraphData{
						Quests: map[string]*dm.Quest{
							x: {Attempts: map[uint32]*dm.Attempt{2: {}, 3: {}}},
						},
					})
				})

				Convey("oob range", func() {
					req.Queries[0] = dm.AttemptRangeQuery(x, 2, 6)
					So(req, WalkShouldReturn(c), &dm.GraphData{
						Quests: map[string]*dm.Quest{
							x: {Attempts: map[uint32]*dm.Attempt{
								2: {}, 3: {}, 4: {}, 5: {DNE: true}}},
						},
					})
				})
			})

			Convey("filtered attempt results", func() {
				x := ensureQuest(c, "x")
				execute(c, dm.NewAttemptID(w, 1))
				depOn(c, dm.NewAttemptID(w, 1), dm.NewAttemptID(x, 1))
				_, err := s.EnsureAttempt(c, &dm.EnsureAttemptReq{ToEnsure: dm.NewAttemptID(x, 2)})
				So(err, ShouldBeNil)
				ttest.Drain(c)

				origNow := datastore.RoundTime(clock.Now(c))

				execute(c, dm.NewAttemptID(x, 1))
				execute(c, dm.NewAttemptID(x, 2))

				exp := google_pb.NewTimestamp(datastore.RoundTime(clock.Now(c).Add(time.Hour * 24 * 4)))

				_, err = s.FinishAttempt(c, &dm.FinishAttemptReq{
					Auth: &dm.Execution_Auth{
						Id:    dm.NewExecutionID(x, 1, 1),
						Token: []byte("sekret"),
					},
					JsonResult: `{"data": ["I can see this"]}`,
					Expiration: exp,
				})
				So(err, ShouldBeNil)

				_, err = s.FinishAttempt(c, &dm.FinishAttemptReq{
					Auth: &dm.Execution_Auth{
						Id:    dm.NewExecutionID(x, 2, 1),
						Token: []byte("sekret"),
					},
					JsonResult: `{"data": ["nope"]}`,
					Expiration: exp,
				})
				So(err, ShouldBeNil)
				ttest.Drain(c)

				execute(c, dm.NewAttemptID(w, 1))

				req.Auth = &dm.Execution_Auth{
					Id:    dm.NewExecutionID(w, 1, 2),
					Token: []byte("sekret"),
				}
				req.Limit.MaxDepth = 2
				req.Include.AttemptResult = true
				req.Queries[0] = dm.AttemptListQueryL(map[string][]uint32{x: nil})

				x1Expect := dm.NewAttemptFinished(exp.Time(), `{"data":["I can see this"]}`)
				x1Expect.Data.Created = google_pb.NewTimestamp(origNow.Add(-time.Second * 12))
				x1Expect.Data.Modified = google_pb.NewTimestamp(origNow)
				x1Expect.Data.NumExecutions = 1

				x2Expect := dm.NewAttemptFinished(exp.Time(), "")
				x2Expect.Data.Created = google_pb.NewTimestamp(origNow.Add(-time.Second * 18))
				x2Expect.Data.Modified = google_pb.NewTimestamp(origNow)
				x2Expect.Data.NumExecutions = 1

				So(req, WalkShouldReturn(c), &dm.GraphData{
					Quests: map[string]*dm.Quest{
						x: {Attempts: map[uint32]*dm.Attempt{
							1: x1Expect,
							2: x2Expect,
						}},
					}})
			})

			Convey("deps (no dest attempts)", func() {
				req.Limit.MaxDepth = 2
				execute(c, dm.NewAttemptID(w, 1))
				x := ensureQuest(c, "x")
				depOn(c, dm.NewAttemptID(w, 1), dm.NewAttemptID(x, 1), dm.NewAttemptID(x, 2))

				Convey("before tumble", func() {
					Convey("deps", func() {
						req.Include.FwdDeps = true
						// didn't run tumble, so that x|1 and x|2 don't get created.
						So(req, WalkShouldReturn(c), &dm.GraphData{
							Quests: map[string]*dm.Quest{
								w: {Attempts: map[uint32]*dm.Attempt{1: {
									FwdDeps: dm.NewAttemptFanout(map[string][]uint32{
										x: {2, 1},
									}),
								}}},
								x: {Attempts: map[uint32]*dm.Attempt{1: {DNE: true}, 2: {DNE: true}}},
							},
						})
					})
				})

				Convey("after tumble", func() {
					ttest.Drain(c)
					Convey("deps (with dest attempts)", func() {
						req.Include.FwdDeps = true
						req.Include.BackDeps = true
						So(req, WalkShouldReturn(c), &dm.GraphData{
							Quests: map[string]*dm.Quest{
								w: {Attempts: map[uint32]*dm.Attempt{1: {
									FwdDeps:  dm.NewAttemptFanout(map[string][]uint32{x: {2, 1}}),
									BackDeps: &dm.AttemptFanout{},
								}}},
								x: {Attempts: map[uint32]*dm.Attempt{1: {
									FwdDeps:  &dm.AttemptFanout{},
									BackDeps: dm.NewAttemptFanout(map[string][]uint32{w: {1}}),
								}, 2: {
									FwdDeps:  &dm.AttemptFanout{},
									BackDeps: dm.NewAttemptFanout(map[string][]uint32{w: {1}}),
								}}},
							},
						})
					})

					Convey("diamond", func() {
						z := ensureQuest(c, "z")
						execute(c, dm.NewAttemptID(x, 1))
						execute(c, dm.NewAttemptID(x, 2))
						depOn(c, dm.NewAttemptID(x, 1), dm.NewAttemptID(z, 1))
						depOn(c, dm.NewAttemptID(x, 2), dm.NewAttemptID(z, 1))
						ttest.Drain(c)

						So(req, WalkShouldReturn(c), &dm.GraphData{
							Quests: map[string]*dm.Quest{
								w: {Attempts: map[uint32]*dm.Attempt{1: {}}},
								x: {Attempts: map[uint32]*dm.Attempt{1: {}, 2: {}}},
								z: {Attempts: map[uint32]*dm.Attempt{1: {}}},
							},
						})
					})

					Convey("diamond (dfs)", func() {
						z := ensureQuest(c, "z")
						execute(c, dm.NewAttemptID(x, 1))
						execute(c, dm.NewAttemptID(x, 2))
						depOn(c, dm.NewAttemptID(x, 1), dm.NewAttemptID(z, 1))
						depOn(c, dm.NewAttemptID(x, 2), dm.NewAttemptID(z, 1))
						ttest.Drain(c)

						req.Mode.Dfs = true
						So(req, WalkShouldReturn(c), &dm.GraphData{
							Quests: map[string]*dm.Quest{
								w: {Attempts: map[uint32]*dm.Attempt{1: {}}},
								x: {Attempts: map[uint32]*dm.Attempt{1: {}, 2: {}}},
								z: {Attempts: map[uint32]*dm.Attempt{1: {}}},
							},
						})
					})

					Convey("attemptlist", func() {
						req.Limit.MaxDepth = 0
						req.Include.ObjectIds = true
						req.Queries[0] = dm.AttemptListQueryL(map[string][]uint32{x: nil})
						So(req, WalkShouldReturn(c), &dm.GraphData{
							Quests: map[string]*dm.Quest{
								x: {
									Id: dm.NewQuestID(x),
									Attempts: map[uint32]*dm.Attempt{
										1: {Id: dm.NewAttemptID(x, 1)},
										2: {Id: dm.NewAttemptID(x, 2)},
									},
								},
							},
						})
					})

				})

				Convey("early stop", func() {
					req.Limit.MaxDepth = 100
					req.Limit.MaxTime = google_pb.NewDuration(time.Nanosecond)
					tc := clock.Get(c).(testclock.TestClock)
					tc.SetTimerCallback(func(d time.Duration, t clock.Timer) {
						tc.Add(d + time.Second)
					})
					ret, err := (&deps{}).WalkGraph(c, req)
					So(err, ShouldBeNil)
					So(ret.HadMore, ShouldBeTrue)
				})

			})
		})
	})
}