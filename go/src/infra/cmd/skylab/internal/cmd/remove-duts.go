// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/maruel/subcommands"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/grpc/prpc"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/cmd/skylab/internal/site"
)

// RemoveDuts subcommand: RemoveDuts a DUT from a drone.
var RemoveDuts = &subcommands.Command{
	UsageLine: "remove-duts [-drone DRONE] [DUT_ID...]",
	ShortDesc: "Remove a DUT from a drone",
	LongDesc:  "Remove a DUT from a drone",
	CommandRun: func() subcommands.CommandRun {
		c := &removeDutsRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)
		c.Flags.StringVar(&c.server, "drone", "", "Drone to remove DUTs from (optional; if absent, remove from current drone)")
		return c
	},
}

type removeDutsRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags
	envFlags  envFlags
	server    string
}

func (c *removeDutsRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", progName, err)
		return 1
	}
	return 0
}

func (c *removeDutsRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	if c.Flags.NArg() == 0 {
		return errors.New("must specify at least 1 DUT")
	}

	req := &fleet.RemoveDutsFromDronesRequest{
		Removals: make([]*fleet.RemoveDutsFromDronesRequest_Item, c.Flags.NArg()),
	}
	for i, dut := range c.Flags.Args() {
		req.Removals[i] = &fleet.RemoveDutsFromDronesRequest_Item{DutId: dut, DroneHostname: c.server}
	}

	ctx := cli.GetContext(a, c, env)
	hc, err := httpClient(ctx, &c.authFlags)
	if err != nil {
		return err
	}
	e := c.envFlags.Env()
	ic := fleet.NewInventoryPRPCClient(&prpc.Client{
		C:       hc,
		Host:    e.AdminService,
		Options: site.DefaultPRPCOptions,
	})

	resp, err := ic.RemoveDutsFromDrones(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.Removed) == 0 {
		fmt.Fprintln(a.GetErr(), "No DUTs removed")
		return nil
	}

	t := tabwriter.NewWriter(a.GetOut(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(t, resp.Url)
	fmt.Fprintln(t, "DUT ID\tRemoved from drone\t")
	for _, r := range resp.Removed {
		fmt.Fprintf(t, "%s\t%s\t\n", r.DutId, r.DroneHostname)
	}
	t.Flush()

	return nil
}