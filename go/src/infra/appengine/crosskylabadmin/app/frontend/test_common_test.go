// Copyright 2018 The LUCI Authors.
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

package frontend

import (
	"testing"

	"infra/appengine/crosskylabadmin/app/clients"

	. "github.com/smartystreets/goconvey/convey"
	swarming "go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/data/strpair"
)

func TestBotContainsDims(t *testing.T) {
	Convey("single dimension matches", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "a", Value: []string{"b"}},
					},
				},
				strpair.Map{"a": []string{"b"}}),
			ShouldBeTrue,
		)
	})

	Convey("requested values are subset of available values", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "a", Value: []string{"b", "c"}},
					},
				},
				strpair.Map{"a": []string{"b"}}),
			ShouldBeTrue,
		)
	})

	Convey("all values match", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "a", Value: []string{"b", "c"}},
					},
				},
				strpair.Map{"a": []string{"b", "c"}}),
			ShouldBeTrue,
		)
	})

	Convey("requested key subset of available keys and value matches", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "a", Value: []string{"b"}},
						{Key: "j", Value: []string{"l"}},
					},
				},
				strpair.Map{"a": []string{"b"}}),
			ShouldBeTrue,
		)
	})

	Convey("keys do not match", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "j", Value: []string{"b"}},
					},
				},
				strpair.Map{"a": []string{"b"}}),
			ShouldBeFalse,
		)
	})

	Convey("values do not match", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "a", Value: []string{"b"}},
					},
				},
				strpair.Map{"a": []string{"c"}}),
			ShouldBeFalse,
		)
	})

	Convey("some values do not match", t, func() {
		So(
			botContainsDims(
				&swarming.SwarmingRpcsBotInfo{
					Dimensions: []*swarming.SwarmingRpcsStringListPair{
						{Key: "a", Value: []string{"b", "d"}},
					},
				},
				strpair.Map{"a": []string{"b", "c"}}),
			ShouldBeFalse,
		)
	})
}

func TestBotForDUT(t *testing.T) {
	Convey("empty dimensions", t, func() {
		So(botForDUT("dut1", "", ""), ShouldResemble, &swarming.SwarmingRpcsBotInfo{
			BotId: "bot_dut1",
			Dimensions: []*swarming.SwarmingRpcsStringListPair{
				{Key: "dut_state", Value: []string{""}},
				{Key: "dut_id", Value: []string{"dut1"}},
				{Key: "dut_name", Value: []string{"dut1-host"}},
			},
		})
	})

	Convey("non-trivial dimensions with whitespace", t, func() {
		So(botForDUT("dut1", "fake_state", "a: x, y ; b :z"), ShouldResemble, &swarming.SwarmingRpcsBotInfo{
			BotId: "bot_dut1",
			Dimensions: []*swarming.SwarmingRpcsStringListPair{
				{Key: "a", Value: []string{"x", "y"}},
				{Key: "b", Value: []string{"z"}},
				{Key: "dut_state", Value: []string{"fake_state"}},
				{Key: "dut_id", Value: []string{"dut1"}},
				{Key: "dut_name", Value: []string{"dut1-host"}},
			},
		},
		)
	})
}

func TestCreateTaskArgsMatcher(t *testing.T) {
	Convey("matching with nil args fails", t, func() {
		So((&createTaskArgsMatcher{
			DutID:        "dut_id",
			DutState:     "dut_state",
			Priority:     10,
			CmdSubString: "-b c",
		}).Matches(nil), ShouldBeFalse)
	})

	Convey("with non-nil args", t, func() {
		arg := &clients.SwarmingCreateTaskArgs{
			Cmd:                  []string{"a", "-b", "c"},
			DutID:                "dut_id",
			DutState:             "dut_state",
			ExecutionTimeoutSecs: 20,
			Priority:             10,
		}
		Convey("all fields matching pass", func() {
			So((&createTaskArgsMatcher{
				DutID:        "dut_id",
				DutState:     "dut_state",
				Priority:     10,
				CmdSubString: "-b c",
			}).Matches(arg), ShouldBeTrue)
		})
		Convey("missing fields matching pass", func() {
			So((&createTaskArgsMatcher{}).Matches(arg), ShouldBeTrue)
		})

		Convey("mistaching DutID fails", func() {
			So((&createTaskArgsMatcher{
				DutID:        "wrong",
				DutState:     "dut_state",
				Priority:     10,
				CmdSubString: "-b c",
			}).Matches(arg), ShouldBeFalse)
		})
		Convey("mistaching DutState fails", func() {
			So((&createTaskArgsMatcher{
				DutID:        "dut_id",
				DutState:     "wrong",
				Priority:     10,
				CmdSubString: "-b c",
			}).Matches(arg), ShouldBeFalse)
		})
		Convey("mistaching Priority fails", func() {
			So((&createTaskArgsMatcher{
				DutID:        "dut_id",
				DutState:     "dut_state",
				Priority:     999,
				CmdSubString: "-b c",
			}).Matches(arg), ShouldBeFalse)
		})
		Convey("mistaching CmdSubstring fails", func() {
			So((&createTaskArgsMatcher{
				DutID:        "dut_id",
				DutState:     "dut_state",
				Priority:     10,
				CmdSubString: "-x z",
			}).Matches(arg), ShouldBeFalse)
		})
	})
}
