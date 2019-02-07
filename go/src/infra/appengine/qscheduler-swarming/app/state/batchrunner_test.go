// Copyright 2019 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/appengine/gaetesting"

	"infra/appengine/qscheduler-swarming/app/state"
	"infra/appengine/qscheduler-swarming/app/state/types"
	"infra/qscheduler/qslib/scheduler"
)

func TestBatcherErrors(t *testing.T) {
	Convey("Given a testing context with a scheduler pool, and a batcher for that pool", t, func() {
		ctx := gaetesting.TestingContext()
		poolID := "pool 1"
		store := state.NewStore(poolID)
		s := types.NewQScheduler(poolID, time.Now())
		store.Save(ctx, s)

		batcher := state.NewBatchRunnerForTest()
		batcher.Start(store)
		defer batcher.Close()

		Convey("an error in one operation in the batch is propagated correctly.", func() {
			var goodError error
			goodOperation := func(ctx context.Context, s *types.QScheduler, m scheduler.MetricsSink) error {
				return nil
			}

			var badError error
			badOperation := func(ctx context.Context, s *types.QScheduler, m scheduler.MetricsSink) error {
				return errors.New("a bad error occurred")
			}

			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				done := batcher.RunOperation(ctx, goodOperation, 0)
				goodError = <-done
				wg.Done()
			}()

			go func() {
				done := batcher.RunOperation(ctx, badOperation, 0)
				badError = <-done
				wg.Done()
			}()

			batcher.TBatchStart()
			batcher.TBatchWait(2)
			wg.Wait()

			So(badError, ShouldNotBeNil)
			So(badError.Error(), ShouldEqual, "a bad error occurred")

			So(goodError, ShouldNotBeNil)
			So(goodError.Error(), ShouldEqual, "batch failed due to error on different request")
		})
	})
}

func TestBatcherBehavior(t *testing.T) {
	Convey("Given a testing context with a scheduler pool, and a batcher for that pool", t, func() {
		ctx := gaetesting.TestingContext()
		poolID := "pool 1"
		store := state.NewStore(poolID)
		s := types.NewQScheduler(poolID, time.Now())
		store.Save(ctx, s)

		batcher := state.NewBatchRunnerForTest()
		batcher.Start(store)
		defer batcher.Close()

		Convey("a batch of requests are run in priority order.", func() {
			s := &[]string{}
			operationA := func(_ context.Context, _ *types.QScheduler, _ scheduler.MetricsSink) error {
				temp := append(*s, "A")
				s = &temp
				return nil
			}
			operationB := func(_ context.Context, _ *types.QScheduler, _ scheduler.MetricsSink) error {
				temp := append(*s, "B")
				s = &temp
				return nil
			}

			wg := sync.WaitGroup{}
			for i := 0; i < 10; i++ {
				wg.Add(2)
				go func() {
					done := batcher.RunOperation(ctx, operationA, state.BatchPriorityNotify)
					<-done
					wg.Done()
				}()
				go func() {
					done := batcher.RunOperation(ctx, operationB, state.BatchPriorityAssign)
					<-done
					wg.Done()
				}()
			}
			batcher.TBatchStart()
			batcher.TBatchWait(20)
			wg.Wait()

			j := strings.Join(*s, "")
			So(j, ShouldEqual, "AAAAAAAAAABBBBBBBBBB")
		})
	})
}