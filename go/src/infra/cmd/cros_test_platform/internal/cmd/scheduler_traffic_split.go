// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"fmt"

	"go.chromium.org/chromiumos/infra/proto/go/test_platform"

	"infra/cmd/cros_test_platform/internal/site"

	"github.com/maruel/subcommands"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/migration/scheduler"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/errors"
)

// SchedulerTrafficSplit implements the `scheduler-traffic-split` subcommand.
var SchedulerTrafficSplit = &subcommands.Command{
	UsageLine: "scheduler-traffic-split -input_json /path/to/input.json -output_json /path/to/output.json",
	ShortDesc: "Determine traffic split between backend schedulers.",
	LongDesc: `Determine traffic split between backend schedulers, i.e. Autotest vs Skylab.

Step input and output is JSON encoded protobuf defined at
https://chromium.googlesource.com/chromiumos/infra/proto/+/master/src/test_platform/steps/scheduler_traffic_split.proto`,
	CommandRun: func() subcommands.CommandRun {
		c := &schedulerTrafficSplitRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.Flags.StringVar(&c.inputPath, "input_json", "", "Path that contains JSON encoded test_platform.steps.SchedulerTrafficSplitRequest")
		c.Flags.StringVar(&c.outputPath, "output_json", "", "Path where JSON encoded test_platform.steps.SchedulerTrafficSplitResponse should be written.")
		return c
	},
}

type schedulerTrafficSplitRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags

	inputPath  string
	outputPath string
}

func (c *schedulerTrafficSplitRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	err := c.innerRun(a, args, env)
	if err != nil {
		fmt.Fprintf(a.GetErr(), "%s\n", err)
	}
	return exitCode(err)
}

func (c *schedulerTrafficSplitRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	if err := c.processCLIArgs(args); err != nil {
		return err
	}
	var request steps.SchedulerTrafficSplitRequest
	if err := readRequest(c.inputPath, &request); err != nil {
		return err
	}
	return errors.Reason("not implemented").Err()
}

func (c *schedulerTrafficSplitRun) processCLIArgs(args []string) error {
	if len(args) > 0 {
		return errors.Reason("have %d positional args, want 0", len(args)).Err()
	}
	if c.inputPath == "" {
		return errors.Reason("-input_json not specified").Err()
	}
	if c.outputPath == "" {
		return errors.Reason("-output_json not specified").Err()
	}
	return nil
}

func determineTrafficSplit(request *steps.SchedulerTrafficSplitRequest, trafficSplitConfig *scheduler.TrafficSplit) (*steps.SchedulerTrafficSplitResponse, error) {
	if err := ensureSufficientForTrafficSplit(request.Request); err != nil {
		return nil, errors.Annotate(err, "determine traffic split").Err()
	}

	rules := determineRelevantRules(request.Request, trafficSplitConfig.Rules)
	switch {
	case len(rules) == 0:
		return nil, errors.Reason("no matching traffic split rule").Err()
	case len(rules) > 1:
		return nil, errors.Reason("too many matching traffic split rules %s for request %s", rules, request).Err()
	default:
		// good case fallthrough.
	}

	r := rules[0]
	switch r.Backend {
	case scheduler.Backend_BACKEND_AUTOTEST:
		return &steps.SchedulerTrafficSplitResponse{
			AutotestRequest: request.Request,
		}, nil
	case scheduler.Backend_BACKEND_SKYLAB:
		return &steps.SchedulerTrafficSplitResponse{
			SkylabRequest: request.Request,
		}, nil
	default:
		return nil, errors.Reason("invalid backend %s in rule", r.Backend.String()).Err()
	}
}

func ensureSufficientForTrafficSplit(r *test_platform.Request) error {
	if r.GetParams().GetScheduling().GetPool() == nil {
		return errors.Reason("request contains no pool information").Err()
	}
	return nil
}

func determineRelevantRules(request *test_platform.Request, rules []*scheduler.Rule) []*scheduler.Rule {
	ret := []*scheduler.Rule{}
	for _, r := range rules {
		if isRuleRelevant(request, r) {
			ret = append(ret, r)
		}
	}
	return ret
}

func isRuleRelevant(request *test_platform.Request, rule *scheduler.Rule) bool {
	if isNonEmptyAndDistinct(
		request.GetParams().GetSoftwareAttributes().GetBuildTarget().GetName(),
		rule.GetRequest().GetBuildTarget().GetName(),
	) {
		return false
	}
	if isNonEmptyAndDistinct(
		request.GetParams().GetHardwareAttributes().GetModel(),
		rule.GetRequest().GetModel(),
	) {
		return false
	}
	return isSchedulingRelevant(request.GetParams().GetScheduling(), rule.GetRequest().GetScheduling())
}

func isSchedulingRelevant(got, want *test_platform.Request_Params_Scheduling) bool {
	if isNonEmptyAndDistinct(got.GetUnmanagedPool(), want.GetUnmanagedPool()) {
		return false
	}
	if isNonEmptyAndDistinct(got.GetManagedPool().String(), want.GetManagedPool().String()) {
		return false
	}
	if isNonEmptyAndDistinct(got.GetQuotaAccount(), want.GetQuotaAccount()) {
		return false
	}
	return true
}

func isNonEmptyAndDistinct(got, want string) bool {
	return got != "" && got != want
}