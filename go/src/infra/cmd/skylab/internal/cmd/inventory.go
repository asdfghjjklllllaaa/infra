// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/maruel/subcommands"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/grpc/prpc"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/cmd/skylab/internal/site"
)

// Inventory subcommand: Print host inventory.
var Inventory = &subcommands.Command{
	UsageLine: "inventory [-dev]",
	ShortDesc: "Print host inventory",
	LongDesc:  "Print host inventory.",
	CommandRun: func() subcommands.CommandRun {
		c := &inventoryRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)
		return c
	},
}

type inventoryRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags
	envFlags  envFlags
}

func (c *inventoryRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", progName, err)
		return 1
	}
	return 0
}

func (c *inventoryRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	ctx := cli.GetContext(a, c, env)
	hc, err := httpClient(ctx, &c.authFlags)
	if err != nil {
		return err
	}
	e := c.envFlags.Env()
	tc := fleet.NewTrackerPRPCClient(&prpc.Client{
		C:       hc,
		Host:    e.AdminService,
		Options: site.DefaultPRPCOptions,
	})
	res, err := tc.SummarizeBots(ctx, &fleet.SummarizeBotsRequest{})
	if err != nil {
		return err
	}
	r := compileInventoryReport(res.GetBots())
	_ = printInventory(os.Stdout, r)
	return nil
}

// inventoryReport contains the compiled status of the inventory
// state.
type inventoryReport struct {
	labs   []*inventoryCount
	models []*inventoryCount
}

// inventoryCount contains the inventory count for some subset of the
// inventory.
type inventoryCount struct {
	name  string
	good  int
	bad   int
	spare int
}

func (ic inventoryCount) available() int {
	return ic.spare - ic.bad
}

func (ic inventoryCount) total() int {
	return ic.good + ic.bad
}

type inventoryMap struct {
	// This is wrapped in a struct to prevent accidentally
	// bypassing the get method.
	m map[string]*inventoryCount
}

func newInventoryMap() inventoryMap {
	return inventoryMap{
		m: make(map[string]*inventoryCount),
	}
}

// Get gets the inventoryCount from the map, allocating a new
// object if needed.
func (m inventoryMap) get(key string) *inventoryCount {
	if ic := m.m[key]; ic != nil {
		return ic
	}
	ic := &inventoryCount{name: key}
	m.m[key] = ic
	return ic
}

// Slice returns a sorted slice of the map's contents.
func (m inventoryMap) slice() []*inventoryCount {
	s := make([]*inventoryCount, 0, len(m.m))
	for _, ic := range m.m {
		s = append(s, ic)
	}
	sort.Slice(s, func(i, j int) bool { return s[i].name < s[j].name })
	return s
}

func compileInventoryReport(bs []*fleet.BotSummary) *inventoryReport {
	labCounts := newInventoryMap()
	modelCounts := newInventoryMap()
	for _, b := range bs {
		d := b.GetDimensions()
		mc := modelCounts.get(d.GetModel())
		addBotCount(mc, b)
		lc := labCounts.get(botLocation(b))
		addBotCount(lc, b)
	}
	return &inventoryReport{
		labs:   labCounts.slice(),
		models: modelCounts.slice(),
	}
}

func botLocation(b *fleet.BotSummary) string {
	n := b.GetDimensions().GetDutName()
	if i := strings.IndexByte(n, '-'); i > -1 {
		return n[:i]
	}
	return n
}

func addBotCount(ic *inventoryCount, b *fleet.BotSummary) {
	if b.Health == fleet.Health_Healthy {
		ic.good++
	} else {
		ic.bad++
	}
	if isSuites(b) {
		ic.spare++
	}
}

// isSuites returns true if the bot is in the suites pool
func isSuites(b *fleet.BotSummary) bool {
	ps := b.GetDimensions().GetPools()
	for _, p := range ps {
		if p == "suites" {
			return true
		}
	}
	return false
}

func printInventory(w io.Writer, r *inventoryReport) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Inventory by location")
	fmt.Fprintln(tw, "===============================================================================")
	fmt.Fprintf(tw, "Location\tAvail\tGood\tBad\tSpare\tTotal\t\n")
	for _, i := range r.labs {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%d\t\n",
			i.name, i.available(), i.good, i.bad, i.spare, i.total())
	}
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "Inventory by model")
	fmt.Fprintln(tw, "===============================================================================")
	fmt.Fprintf(tw, "Model\tAvail\tGood\tBad\tSpare\tTotal\t\n")
	for _, i := range r.models {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%d\t\n",
			i.name, i.available(), i.good, i.bad, i.spare, i.total())
	}
	return tw.Flush()
}