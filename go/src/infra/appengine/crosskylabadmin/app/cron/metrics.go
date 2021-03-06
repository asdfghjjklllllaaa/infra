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

package cron

import (
	"go.chromium.org/luci/common/tsmon/field"
	"go.chromium.org/luci/common/tsmon/metric"
)

var (
	refreshBotsTick = metric.NewCounter(
		"chromeos/crosskylabadmin/cron/refresh_bots",
		"RefreshBots cron attempt",
		nil,
		field.Bool("success"),
	)
	ensureBackgroundTasksTick = metric.NewCounter(
		"chromeos/crosskylabadmin/cron/ensure_background_tasks",
		"EnsureBackgroundTasks cron attempt",
		nil,
		field.Bool("success"),
	)
	triggerRepairOnIdleTick = metric.NewCounter(
		"chromeos/crosskylabadmin/cron/trigger_repair_on_idle",
		"TriggerRepairOnIdle cron attempt",
		nil,
		field.Bool("success"),
	)
	triggerRepairOnRepairFailedTick = metric.NewCounter(
		"chromeos/crosskylabadmin/cron/trigger_repair_on_repair_failed",
		"TriggerRepairOnRepairFailed cron attempt",
		nil,
		field.Bool("success"),
	)
)
