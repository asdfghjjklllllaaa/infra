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

// Package botsummary implements datastore bot summary access.
package botsummary

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/luci/common/errors"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
)

// botSummaryKind is the datastore entity kind for fleetBotSummaryEntity.
const botSummaryKind = "fleetBotSummary"

// Entity is a datastore entity that stores fleet.BotSummary in
// protobuf binary format.
//
// In effect, this is a cache of task and bot history information
// from the Swarming server over a fixed time period.
type Entity struct {
	_kind string `gae:"$kind,fleetBotSummary"`
	// Swarming bot's dut_id dimension.
	//
	// This dimension is an opaque reference to the managed DUT's uuid in skylab
	// inventory data.
	DutID string `gae:"$id"`
	// BotID is the unique ID of the swarming bot.
	BotID string
	// Data is the fleet.BotSummary object serialized to protobuf binary format.
	Data []byte `gae:",noindex"`
}

// Decode converts the Entity into a fleet.BotSummary.
func (e *Entity) Decode() (*fleet.BotSummary, error) {
	var bs fleet.BotSummary
	if err := proto.Unmarshal(e.Data, &bs); err != nil {
		return nil, errors.Annotate(err, "failed to unmarshal bot summary for bot with dut_id %q", e.DutID).Err()
	}
	return &bs, nil
}

// Insert inserts fleet.BotSummary into datastore.
func Insert(ctx context.Context, bsm map[string]*fleet.BotSummary) (dutIDs []string, err error) {
	updated := make([]string, 0, len(bsm))
	es := make([]*Entity, 0, len(bsm))
	for bid, bs := range bsm {
		data, err := proto.Marshal(bs)
		if err != nil {
			return nil, errors.Annotate(err, "failed to marshal BotSummary for dut %q", bs.DutId).Err()
		}
		es = append(es, &Entity{
			DutID: bs.DutId,
			BotID: bid,
			Data:  data,
		})
		updated = append(updated, bs.DutId)
	}
	if err := datastore.Put(ctx, es); err != nil {
		return nil, errors.Annotate(err, "failed to put BotSummaries").Err()
	}
	return updated, nil
}

// Get gets Entities from datastore.
func Get(ctx context.Context, sels []*fleet.BotSelector) ([]*Entity, error) {
	// No selectors implies summarize all bots.
	if len(sels) == 0 {
		return GetAll(ctx)
	}

	dutIDs := make([]string, 0, len(sels))
	for _, s := range sels {
		// datastore rejects search for empty key with InvalidKey.
		// For us, this is simply an impossible filter.
		if s.DutId == "" {
			continue
		}
		dutIDs = append(dutIDs, s.DutId)
	}
	return GetByDutID(ctx, dutIDs)
}

// GetAll gets all Entities from the datastore.
func GetAll(ctx context.Context) ([]*Entity, error) {
	var es []*Entity
	q := datastore.NewQuery(botSummaryKind)
	err := datastore.GetAll(ctx, q, &es)
	if err != nil {
		return nil, errors.Annotate(err, "failed to get all bots from datastore").Err()
	}
	return es, nil
}

// GetByDutID gets Entities from datastore by DUT ID.  Missing DUT IDs
// are ignored.  Successfully fetched Entities are returned even if
// others encountered errors.
func GetByDutID(ctx context.Context, dutIDs []string) ([]*Entity, error) {
	es := make([]*Entity, 0, len(dutIDs))
	for _, id := range dutIDs {
		es = append(es, &Entity{DutID: id})
	}
	switch err := datastore.Get(ctx, es); err := err.(type) {
	case nil:
		return es, nil
	case errors.MultiError:
		if len(es) != len(err) {
			panic(fmt.Sprintf("bot summary length %d != multierror %d",
				len(es), len(err)))
		}
		es = removeErroredEntities(es, err)
		if datastore.IsErrNoSuchEntity(err) {
			return es, nil
		}
		return es, err
	default:
		return nil, err
	}
}

// removeErroredEntities returns a slice of Entities without the ones
// with errors.
func removeErroredEntities(es []*Entity, merr errors.MultiError) []*Entity {
	ok := make([]*Entity, 0, len(es))
	for i, e := range es {
		if merr[i] == nil {
			ok = append(ok, e)
		}
	}
	return ok
}