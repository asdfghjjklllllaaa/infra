// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/maruel/subcommands"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"

	qscheduler "infra/appengine/qscheduler-swarming/api/qscheduler/v1"
	"infra/cmd/qscheduler/internal/site"
	"infra/qscheduler/qslib/scheduler"
)

// AddAccount subcommand: add an account.
var AddAccount = &subcommands.Command{
	UsageLine: "add-account",
	ShortDesc: "Add a quota account",
	LongDesc:  "Add a quota account",
	CommandRun: func() subcommands.CommandRun {
		c := &addAccountRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)
		c.Flags.StringVar(&c.poolID, "id", "", "Scheduler ID to modify.")
		c.Flags.StringVar(&c.accountID, "account", "", "Account ID to create.")
		c.Flags.Var(MultiFloat(&c.chargeRates), "rate", "Quota recharge rate for a given priority level. "+
			"May be specified multiple times, to specify charge rate at P0, P1, P2, ...")
		c.Flags.Float64Var(&c.chargeTime, "charge-time", 0,
			"Maximum amount of time (seconds) for which the account can accumulate quota.")
		c.Flags.IntVar(&c.fanout, "fanout", 0, "Maximum number of concurrent tasks that account will pay for.")
		return c
	},
}

type addAccountRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags
	envFlags  envFlags

	poolID      string
	accountID   string
	chargeRates []float64
	chargeTime  float64
	fanout      int
}

func (c *addAccountRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if c.poolID == "" {
		fmt.Fprintf(os.Stderr, "Must specify id.\n")
		return 1
	}

	if c.accountID == "" {
		fmt.Fprintf(os.Stderr, "Must specify account id.\n")
		return 1
	}

	ctx := cli.GetContext(a, c, env)

	adminClient, err := newAdminClient(ctx, &c.authFlags, &c.envFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create qsadmin client, due to error: %s\n", err.Error())
		return 1
	}

	req := &qscheduler.CreateAccountRequest{
		AccountId: c.accountID,
		PoolId:    c.poolID,
		Config: &scheduler.AccountConfig{
			ChargeRate:       c.chargeRates,
			MaxChargeSeconds: c.chargeTime,
			MaxFanout:        int32(c.fanout),
		},
	}

	_, err = adminClient.CreateAccount(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add account, due to error: %s\n", err.Error())
		return 1
	}

	fmt.Printf("Added account %s to scheduler %s.\n", c.accountID, c.poolID)
	return 0
}