// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package skylab implements logic necessary for Skylab execution of an
// ExecuteRequest.
package skylab

import (
	"context"
	"time"

	build_api "go.chromium.org/chromiumos/infra/proto/go/chromite/api"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"
	swarming_api "go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/swarming/proto/jsonrpc"

	"infra/cmd/cros_test_platform/internal/execution/isolate"
	"infra/cmd/cros_test_platform/internal/execution/swarming"
	"infra/libs/skylab/inventory"
	"infra/libs/skylab/inventory/autotest/labels"
	"infra/libs/skylab/request"
	"infra/libs/skylab/worker"
)

// TaskSet encapsulates the running state of a set of tasks, to satisfy
// a Skylab Execution.
type TaskSet struct {
	testRuns []*testRun
	params   *test_platform.Request_Params
}

type testRun struct {
	test     *build_api.AutotestTest
	attempts []attempt
}

func (t *testRun) RequestArgs() (request.Args, error) {
	isClient, err := t.isClientTest()
	if err != nil {
		return request.Args{}, errors.Annotate(err, "create request args").Err()
	}

	// TODO(akeshet): Run cmd.Config() with correct environment.
	cmd := &worker.Command{
		TaskName:   t.test.Name,
		ClientTest: isClient,
	}

	args := request.Args{
		Cmd:               *cmd,
		SchedulableLabels: toInventoryLabels(t.test.Dependencies),
		// TODO(akeshet): Determine parent task ID correctly.
		ParentTaskID: "",
		// TODO(akeshet): Determine priority correctly.
		Priority: 0,
		// TODO(akeshet): Determine provisionable dimensions correctly.
		ProvisionableDimensions: nil,
		// TODO(akeshet): Determine tags correctly.
		SwarmingTags: nil,
		// TODO(akeshet): Determine timeout correctly.
		Timeout: 0,
	}

	return args, nil
}

func (t *testRun) isClientTest() (bool, error) {
	isClient, ok := isClientTest[t.test.ExecutionEnvironment]
	if !ok {
		return false, errors.Reason("unknown exec environment %s", t.test.ExecutionEnvironment).Err()
	}
	return isClient, nil
}

type attempt struct {
	taskID    string
	completed bool
	state     jsonrpc.TaskState
}

// NewTaskSet creates a new TaskSet.
func NewTaskSet(tests []*build_api.AutotestTest, params *test_platform.Request_Params) *TaskSet {
	testRuns := make([]*testRun, len(tests))
	for i, test := range tests {
		testRuns[i] = &testRun{test: test}
	}
	return &TaskSet{
		testRuns: testRuns,
		params:   params,
	}
}

// LaunchAndWait launches a skylab execution and waits for it to complete,
// polling for new results periodically (TODO(akeshet): and retrying tests that
// need retry, based on retry policy).
//
// If the supplied context is cancelled prior to completion, or some other error
// is encountered, this method returns whatever partial execution response
// was visible to it prior to that error.
func (r *TaskSet) LaunchAndWait(ctx context.Context, swarming swarming.Client, isolate isolate.Client) error {
	if err := r.launch(ctx, swarming); err != nil {
		return err
	}

	return r.wait(ctx, swarming)
}

var isClientTest = map[build_api.AutotestTest_ExecutionEnvironment]bool{
	build_api.AutotestTest_EXECUTION_ENVIRONMENT_CLIENT: true,
	build_api.AutotestTest_EXECUTION_ENVIRONMENT_SERVER: false,
}

func (r *TaskSet) launch(ctx context.Context, swarming swarming.Client) error {
	for _, testRun := range r.testRuns {
		args, err := testRun.RequestArgs()
		if err != nil {
			return errors.Annotate(err, "launch test named %s", testRun.test.Name).Err()
		}

		req, err := request.New(args)
		if err != nil {
			return errors.Annotate(err, "launch test named %s", testRun.test.Name).Err()
		}

		resp, err := swarming.CreateTask(ctx, req)
		if err != nil {
			return errors.Annotate(err, "launch test named %s", testRun.test.Name).Err()
		}

		testRun.attempts = append(testRun.attempts, attempt{taskID: resp.TaskId})
	}
	return nil
}

func (r *TaskSet) wait(ctx context.Context, swarming swarming.Client) error {
	for {
		complete, err := r.tick(ctx, swarming)
		if complete || err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return errors.Annotate(ctx.Err(), "wait for tests").Err()
		case <-clock.After(ctx, 15*time.Second):
		}
	}
}

func (r *TaskSet) tick(ctx context.Context, client swarming.Client) (complete bool, err error) {
	complete = true

	for _, testRun := range r.testRuns {
		attempt := &testRun.attempts[len(testRun.attempts)-1]
		if attempt.completed {
			continue
		}

		results, err := client.GetResults(ctx, []string{attempt.taskID})
		if err != nil {
			return false, errors.Annotate(err, "wait for tests").Err()
		}

		result, err := unpackResult(results, attempt.taskID)
		if err != nil {
			return false, errors.Annotate(err, "wait for tests").Err()
		}

		state, err := swarming.AsTaskState(result.State)
		if err != nil {
			return false, errors.Annotate(err, "wait for tests").Err()
		}
		attempt.state = state

		if !swarming.UnfinishedTaskStates[state] {
			attempt.completed = true
			continue
		}

		// At least one task is not complete.
		complete = false
	}

	return complete, nil
}

func toInventoryLabels(deps []*build_api.AutotestTaskDependency) inventory.SchedulableLabels {
	flatDims := make([]string, len(deps))
	for i, dep := range deps {
		flatDims[i] = dep.Label
	}
	return *labels.Revert(flatDims)
}

func unpackResult(results []*swarming_api.SwarmingRpcsTaskResult, taskID string) (*swarming_api.SwarmingRpcsTaskResult, error) {
	if len(results) != 1 {
		return nil, errors.Reason("expected 1 result for task id %s, got %d", taskID, len(results)).Err()
	}

	result := results[0]
	if result.TaskId != taskID {
		return nil, errors.Reason("expected result for task id %s, got %s", taskID, result.TaskId).Err()
	}

	return result, nil
}

var taskStateToLifeCycle = map[jsonrpc.TaskState]test_platform.TaskState_LifeCycle{
	jsonrpc.TaskState_BOT_DIED:  test_platform.TaskState_LIFE_CYCLE_ABORTED,
	jsonrpc.TaskState_CANCELED:  test_platform.TaskState_LIFE_CYCLE_CANCELLED,
	jsonrpc.TaskState_COMPLETED: test_platform.TaskState_LIFE_CYCLE_COMPLETED,
	// TODO(akeshet): This mapping is inexact. Add a lifecycle entry for this.
	jsonrpc.TaskState_EXPIRED:     test_platform.TaskState_LIFE_CYCLE_CANCELLED,
	jsonrpc.TaskState_KILLED:      test_platform.TaskState_LIFE_CYCLE_ABORTED,
	jsonrpc.TaskState_NO_RESOURCE: test_platform.TaskState_LIFE_CYCLE_REJECTED,
	jsonrpc.TaskState_PENDING:     test_platform.TaskState_LIFE_CYCLE_PENDING,
	jsonrpc.TaskState_RUNNING:     test_platform.TaskState_LIFE_CYCLE_RUNNING,
	// TODO(akeshet): This mapping is inexact. Add a lifecycle entry for this.
	jsonrpc.TaskState_TIMED_OUT: test_platform.TaskState_LIFE_CYCLE_ABORTED,
}

// Response constructs a response based on the current state of the
// TaskSet.
func (r *TaskSet) Response(swarming swarming.URLer) *steps.ExecuteResponse {
	resp := &steps.ExecuteResponse{}
	for _, test := range r.testRuns {
		for _, attempt := range test.attempts {
			resp.TaskResults = append(resp.TaskResults, &steps.ExecuteResponse_TaskResult{
				Name: test.test.Name,
				State: &test_platform.TaskState{
					LifeCycle: taskStateToLifeCycle[attempt.state],
					// TODO(akeshet): Determine a way to extract and identify
					// test verdicts.
					Verdict: test_platform.TaskState_VERDICT_NO_VERDICT,
				},
				TaskId:  attempt.taskID,
				TaskUrl: swarming.GetTaskURL(attempt.taskID),
			})
		}
	}
	return resp
}
