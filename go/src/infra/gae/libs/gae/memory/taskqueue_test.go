// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package memory

import (
	"fmt"
	"infra/gae/libs/gae"
	"math/rand"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"

	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/clock/testclock"
)

func TestTaskQueue(t *testing.T) {
	t.Parallel()

	Convey("TaskQueue", t, func() {
		now := time.Date(2000, time.January, 1, 1, 1, 1, 1, time.UTC)
		c, tc := testclock.UseTime(context.Background(), now)
		c = gae.SetMathRand(c, rand.New(rand.NewSource(clock.Now(c).UnixNano())))
		c = Use(c)

		tq := gae.GetTQ(c).(interface {
			gae.TaskQueue
			gae.TQTestable
		})

		So(tq, ShouldNotBeNil)

		Convey("implements TQMultiReadWriter", func() {
			Convey("Add", func() {
				t := &gae.TQTask{Path: "/hello/world"}

				Convey("works", func() {
					t.Delay = 4 * time.Second
					t.Header = http.Header{}
					t.Header.Add("Cat", "tabby")
					t.Payload = []byte("watwatwat")
					t.RetryOptions = &gae.TQRetryOptions{AgeLimit: 7 * time.Second}
					_, err := tq.Add(t, "")
					So(err, ShouldBeNil)
					name := "Z_UjshxM9ecyMQfGbZmUGOEcgxWU0_5CGLl_-RntudwAw2DqQ5-58bzJiWQN4OKzeuUb9O4JrPkUw2rOvk2Ax46THojnQ6avBQgZdrKcJmrwQ6o4qKfJdiyUbGXvy691yRfzLeQhs6cBhWrgf3wH-VPMcA4SC-zlbJ2U8An7I0zJQA5nBFnMNoMgT-2peGoay3rCSbj4z9VFFm9kS_i6JCaQH518ujLDSNCYdjTq6B6lcWrZAh0U_q3a1S2nXEwrKiw_t9MTNQFgAQZWyGBbvZQPmeRYtu8SPaWzTfd25v_YWgBuVL2rRSPSMvlDwE04nNdtvVzE8vNNiA1zRimmdzKeqATQF9_ReUvj4D7U8dcS703DZWfKMBLgBffY9jqCassOOOw77V72Oq5EVauUw3Qw0L6bBsfM9FtahTKUdabzRZjXUoze3EK4KXPt3-wdidau-8JrVf2XFocjjZbwHoxcGvbtT3b4nGLDlgwdC00bwaFBZWff"
					So(*tq.GetScheduledTasks()["default"][name], ShouldResemble, gae.TQTask{
						ETA:          now.Add(4 * time.Second),
						Header:       http.Header{"Cat": []string{"tabby"}},
						Method:       "POST",
						Name:         name,
						Path:         "/hello/world",
						Payload:      []byte("watwatwat"),
						RetryOptions: &gae.TQRetryOptions{AgeLimit: 7 * time.Second},
					})
				})

				Convey("cannot add to bad queues", func() {
					_, err := tq.Add(nil, "waaat")
					So(err.Error(), ShouldContainSubstring, "UNKNOWN_QUEUE")

					Convey("but you can add Queues when testing", func() {
						tq.CreateQueue("waaat")
						_, err := tq.Add(t, "waaat")
						So(err, ShouldBeNil)

						Convey("you just can't add them twice", func() {
							So(func() { tq.CreateQueue("waaat") }, ShouldPanic)
						})
					})
				})

				Convey("requires a URL", func() {
					t.Path = ""
					tr, err := tq.Add(t, "")
					So(err.Error(), ShouldContainSubstring, "INVALID_URL")
					So(tr, ShouldBeNil)
				})

				Convey("cannot add twice", func() {
					t.Name = "bob"
					_, err := tq.Add(t, "")
					So(err, ShouldBeNil)

					// can't add the same one twice!
					_, err = tq.Add(t, "")
					So(err, ShouldEqual, gae.ErrTQTaskAlreadyAdded)
				})

				Convey("cannot add deleted task", func() {
					t.Name = "bob"
					_, err := tq.Add(t, "")
					So(err, ShouldBeNil)

					err = tq.Delete(t, "")
					So(err, ShouldBeNil)

					// can't add a deleted task!
					_, err = tq.Add(t, "")
					So(err, ShouldEqual, gae.ErrTQTaskAlreadyAdded)
				})

				Convey("cannot set ETA+Delay", func() {
					t.ETA = clock.Now(c).Add(time.Hour)
					tc.Add(time.Second)
					t.Delay = time.Hour
					So(func() { tq.Add(t, "") }, ShouldPanic)
				})

				Convey("must use a reasonable method", func() {
					t.Method = "Crystal"
					_, err := tq.Add(t, "")
					So(err.Error(), ShouldContainSubstring, "bad method")
				})

				Convey("payload gets dumped for non POST/PUT methods", func() {
					t.Method = "HEAD"
					t.Payload = []byte("coool")
					tq, err := tq.Add(t, "")
					So(err, ShouldBeNil)
					So(tq.Payload, ShouldBeNil)

					// check that it didn't modify our original
					So(t.Payload, ShouldResemble, []byte("coool"))
				})

				Convey("invalid names are rejected", func() {
					t.Name = "happy times"
					_, err := tq.Add(t, "")
					So(err.Error(), ShouldContainSubstring, "INVALID_TASK_NAME")
				})

				Convey("can be broken", func() {
					tq.BreakFeatures(nil, "Add")
					_, err := tq.Add(t, "")
					So(err.Error(), ShouldContainSubstring, "TRANSIENT_ERROR")
				})

				Convey("AddMulti also works", func() {
					t2 := dupTask(t)
					t2.Path = "/hi/city"

					expect := []*gae.TQTask{t, t2}

					tasks, err := tq.AddMulti(expect, "default")
					So(err, ShouldBeNil)
					So(len(tasks), ShouldEqual, 2)
					So(len(tq.GetScheduledTasks()["default"]), ShouldEqual, 2)

					for i := range expect {
						Convey(fmt.Sprintf("task %d: %s", i, expect[i].Path), func() {
							expect[i].Method = "POST"
							expect[i].ETA = now
							So(expect[i].Name, ShouldEqual, "")
							So(len(tasks[i].Name), ShouldEqual, 500)
							tasks[i].Name = ""
							So(tasks[i], ShouldResemble, expect[i])
						})
					}

					Convey("can be broken", func() {
						tq.BreakFeatures(nil, "AddMulti")
						_, err := tq.AddMulti([]*gae.TQTask{t}, "")
						So(err.Error(), ShouldContainSubstring, "TRANSIENT_ERROR")
					})

					Convey("is not broken by Add", func() {
						tq.BreakFeatures(nil, "Add")
						_, err := tq.AddMulti([]*gae.TQTask{t}, "")
						So(err, ShouldBeNil)
					})
				})
			})

			Convey("Delete", func() {
				t := &gae.TQTask{Path: "/hello/world"}
				tEnQ, err := tq.Add(t, "")
				So(err, ShouldBeNil)

				Convey("works", func() {
					t.Name = tEnQ.Name
					err := tq.Delete(t, "")
					So(err, ShouldBeNil)
					So(len(tq.GetScheduledTasks()["default"]), ShouldEqual, 0)
					So(len(tq.GetTombstonedTasks()["default"]), ShouldEqual, 1)
					So(tq.GetTombstonedTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
				})

				Convey("cannot delete a task twice", func() {
					err := tq.Delete(tEnQ, "")
					So(err, ShouldBeNil)

					err = tq.Delete(tEnQ, "")
					So(err.Error(), ShouldContainSubstring, "TOMBSTONED_TASK")

					Convey("but you can if you do a reset", func() {
						tq.ResetTasks()

						tEnQ, err := tq.Add(t, "")
						So(err, ShouldBeNil)
						err = tq.Delete(tEnQ, "")
						So(err, ShouldBeNil)
					})
				})

				Convey("cannot delete from bogus queues", func() {
					err := tq.Delete(t, "wat")
					So(err.Error(), ShouldContainSubstring, "UNKNOWN_QUEUE")
				})

				Convey("cannot delete a missing task", func() {
					t.Name = "tarntioarenstyw"
					err := tq.Delete(t, "")
					So(err.Error(), ShouldContainSubstring, "UNKNOWN_TASK")
				})

				Convey("can be broken", func() {
					tq.BreakFeatures(nil, "Delete")
					err := tq.Delete(t, "")
					So(err.Error(), ShouldContainSubstring, "TRANSIENT_ERROR")
				})

				Convey("DeleteMulti also works", func() {
					t2 := dupTask(t)
					t2.Path = "/hi/city"
					tEnQ2, err := tq.Add(t2, "")
					So(err, ShouldBeNil)

					Convey("usually works", func() {
						err = tq.DeleteMulti([]*gae.TQTask{tEnQ, tEnQ2}, "")
						So(err, ShouldBeNil)
						So(len(tq.GetScheduledTasks()["default"]), ShouldEqual, 0)
						So(len(tq.GetTombstonedTasks()["default"]), ShouldEqual, 2)
					})

					Convey("can be broken", func() {
						tq.BreakFeatures(nil, "DeleteMulti")
						err = tq.DeleteMulti([]*gae.TQTask{tEnQ, tEnQ2}, "")
						So(err.Error(), ShouldContainSubstring, "TRANSIENT_ERROR")
					})

					Convey("is not broken by Delete", func() {
						tq.BreakFeatures(nil, "Delete")
						err = tq.DeleteMulti([]*gae.TQTask{tEnQ, tEnQ2}, "")
						So(err, ShouldBeNil)
					})
				})
			})
		})

		Convey("works with transactions", func() {
			t := &gae.TQTask{Path: "/hello/world"}
			tEnQ, err := tq.Add(t, "")
			So(err, ShouldBeNil)

			t2 := &gae.TQTask{Path: "/hi/city"}
			tEnQ2, err := tq.Add(t2, "")
			So(err, ShouldBeNil)

			err = tq.Delete(tEnQ2, "")
			So(err, ShouldBeNil)

			Convey("can view regular tasks", func() {
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					tq := gae.GetTQ(c).(interface {
						gae.TQTestable
						gae.TaskQueue
					})

					So(tq.GetScheduledTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
					So(tq.GetTombstonedTasks()["default"][tEnQ2.Name], ShouldResemble, tEnQ2)
					So(tq.GetTransactionTasks()["default"], ShouldBeNil)
					return nil
				}, nil)
			})

			Convey("can add a new task", func() {
				tEnQ3 := (*gae.TQTask)(nil)

				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					tq := gae.GetTQ(c).(interface {
						gae.TQTestable
						gae.TaskQueue
					})

					t3 := &gae.TQTask{Path: "/sandwitch/victory"}
					tEnQ3, err = tq.Add(t3, "")
					So(err, ShouldBeNil)

					So(tq.GetScheduledTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
					So(tq.GetTombstonedTasks()["default"][tEnQ2.Name], ShouldResemble, tEnQ2)
					So(tq.GetTransactionTasks()["default"][0], ShouldResemble, tEnQ3)
					return nil
				}, nil)

				// name gets generated at transaction-commit-time
				for name := range tq.GetScheduledTasks()["default"] {
					if name == tEnQ.Name {
						continue
					}
					tEnQ3.Name = name
					break
				}

				So(tq.GetScheduledTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
				So(tq.GetScheduledTasks()["default"][tEnQ3.Name], ShouldResemble, tEnQ3)
				So(tq.GetTombstonedTasks()["default"][tEnQ2.Name], ShouldResemble, tEnQ2)
				So(tq.GetTransactionTasks()["default"], ShouldBeNil)
			})

			Convey("can a new task (but reset the state in a test)", func() {
				tEnQ3 := (*gae.TQTask)(nil)

				ttq := interface {
					gae.TQTestable
					gae.TaskQueue
				}(nil)

				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					ttq = gae.GetTQ(c).(interface {
						gae.TQTestable
						gae.TaskQueue
					})

					t3 := &gae.TQTask{Path: "/sandwitch/victory"}
					tEnQ3, err = ttq.Add(t3, "")
					So(err, ShouldBeNil)

					So(ttq.GetScheduledTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
					So(ttq.GetTombstonedTasks()["default"][tEnQ2.Name], ShouldResemble, tEnQ2)
					So(ttq.GetTransactionTasks()["default"][0], ShouldResemble, tEnQ3)

					ttq.ResetTasks()

					So(len(ttq.GetScheduledTasks()["default"]), ShouldEqual, 0)
					So(len(ttq.GetTombstonedTasks()["default"]), ShouldEqual, 0)
					So(len(ttq.GetTransactionTasks()["default"]), ShouldEqual, 0)

					return nil
				}, nil)

				So(len(tq.GetScheduledTasks()["default"]), ShouldEqual, 0)
				So(len(tq.GetTombstonedTasks()["default"]), ShouldEqual, 0)
				So(len(tq.GetTransactionTasks()["default"]), ShouldEqual, 0)

				Convey("and reusing a closed context is bad times", func() {
					_, err := ttq.Add(nil, "")
					So(err.Error(), ShouldContainSubstring, "expired")
				})
			})

			Convey("you can AddMulti as well", func() {
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					tq := gae.GetTQ(c).(interface {
						gae.TQTestable
						gae.TaskQueue
					})
					_, err := tq.AddMulti([]*gae.TQTask{t, t, t}, "")
					So(err, ShouldBeNil)
					So(len(tq.GetScheduledTasks()["default"]), ShouldEqual, 1)
					So(len(tq.GetTransactionTasks()["default"]), ShouldEqual, 3)
					return nil
				}, nil)
				So(len(tq.GetScheduledTasks()["default"]), ShouldEqual, 4)
				So(len(tq.GetTransactionTasks()["default"]), ShouldEqual, 0)
			})

			Convey("unless you add too many things", func() {
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					for i := 0; i < 5; i++ {
						_, err = gae.GetTQ(c).Add(t, "")
						So(err, ShouldBeNil)
					}
					_, err = gae.GetTQ(c).Add(t, "")
					So(err.Error(), ShouldContainSubstring, "BAD_REQUEST")
					return nil
				}, nil)
			})

			Convey("unless you Add to a bad queue", func() {
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					_, err = gae.GetTQ(c).Add(t, "meat")
					So(err.Error(), ShouldContainSubstring, "UNKNOWN_QUEUE")

					Convey("unless you add it!", func() {
						gae.GetTQ(c).(gae.TQTestable).CreateQueue("meat")
						_, err = gae.GetTQ(c).Add(t, "meat")
						So(err, ShouldBeNil)
					})

					return nil
				}, nil)
			})

			Convey("unless Add is broken", func() {
				tq.BreakFeatures(nil, "Add")
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					_, err = gae.GetTQ(c).Add(t, "")
					So(err.Error(), ShouldContainSubstring, "TRANSIENT_ERROR")
					return nil
				}, nil)
			})

			Convey("unless AddMulti is broken", func() {
				tq.BreakFeatures(nil, "AddMulti")
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					_, err = gae.GetTQ(c).AddMulti(nil, "")
					So(err.Error(), ShouldContainSubstring, "TRANSIENT_ERROR")
					return nil
				}, nil)
			})

			Convey("No other features are available, however", func() {
				err := error(nil)
				func() {
					defer func() { err = recover().(error) }()
					gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
						gae.GetTQ(c).Delete(t, "")
						return nil
					}, nil)
				}()
				So(err.Error(), ShouldContainSubstring, "TaskQueue.Delete")
			})

			Convey("adding a new task only happens if we don't errout", func() {
				gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
					t3 := &gae.TQTask{Path: "/sandwitch/victory"}
					_, err = gae.GetTQ(c).Add(t3, "")
					So(err, ShouldBeNil)
					return fmt.Errorf("nooooo")
				}, nil)

				So(tq.GetScheduledTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
				So(tq.GetTombstonedTasks()["default"][tEnQ2.Name], ShouldResemble, tEnQ2)
				So(tq.GetTransactionTasks()["default"], ShouldBeNil)
			})

			Convey("likewise, a panic doesn't schedule anything", func() {
				func() {
					defer func() { recover() }()
					gae.GetRDS(c).RunInTransaction(func(c context.Context) error {
						tq := gae.GetTQ(c).(interface {
							gae.TQTestable
							gae.TaskQueue
						})

						t3 := &gae.TQTask{Path: "/sandwitch/victory"}
						_, err = tq.Add(t3, "")
						So(err, ShouldBeNil)

						panic(fmt.Errorf("nooooo"))
					}, nil)
				}()

				So(tq.GetScheduledTasks()["default"][tEnQ.Name], ShouldResemble, tEnQ)
				So(tq.GetTombstonedTasks()["default"][tEnQ2.Name], ShouldResemble, tEnQ2)
				So(tq.GetTransactionTasks()["default"], ShouldBeNil)
			})

		})
	})
}
