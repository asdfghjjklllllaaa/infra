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

// Package inventory implements the fleet.Inventory service end-points of
// corsskylabadmin.
package inventory

import (
	"fmt"
	"math/rand"
	"sort"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/appengine/crosskylabadmin/app/config"
	"infra/appengine/crosskylabadmin/app/frontend/internal/gitstore"
	"infra/libs/skylab/inventory"

	"github.com/golang/protobuf/proto"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/common/retry"
	"go.chromium.org/luci/grpc/grpcutil"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// Number of drones least loaded drones to consider when picking a drone for a DUT.
	//
	// A random drone from these is assigned to the DUT.
	leastLoadedDronesCount = 8
)

// AssignDutsToDrones implements the method from fleet.InventoryServer interface.
func (is *ServerImpl) AssignDutsToDrones(ctx context.Context, req *fleet.AssignDutsToDronesRequest) (resp *fleet.AssignDutsToDronesResponse, err error) {
	defer func() {
		err = grpcutil.GRPCifyAndLogErr(ctx, err)
	}()
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	s, err := is.newStore(ctx)
	if err != nil {
		return nil, err
	}
	f := func() error {
		if err := s.Refresh(ctx); err != nil {
			return err
		}
		da := newDUTAssigner(ctx, s)
		assigned := make([]*fleet.AssignDutsToDronesResponse_Item, 0, len(req.Assignments))
		for _, a := range req.Assignments {
			i, err := da.assignDUT(ctx, a)
			if err != nil {
				return err
			}
			assigned = append(assigned, i)
		}
		url, err := s.Commit(ctx, "assign DUTs")
		if err != nil {
			return err
		}
		resp = &fleet.AssignDutsToDronesResponse{
			Assigned: assigned,
			Url:      url,
		}
		return nil
	}
	if err = retry.Retry(ctx, transientErrorRetries(), f, retry.LogCallback(ctx, "assignDutsToDronesNoRetry")); err != nil {
		return nil, err
	}
	return resp, nil
}

// RemoveDutsFromDrones implements the method from fleet.InventoryServer interface.
func (is *ServerImpl) RemoveDutsFromDrones(ctx context.Context, req *fleet.RemoveDutsFromDronesRequest) (resp *fleet.RemoveDutsFromDronesResponse, err error) {
	defer func() {
		err = grpcutil.GRPCifyAndLogErr(ctx, err)
	}()
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	s, err := is.newStore(ctx)
	if err != nil {
		return nil, err
	}
	f := func() error {
		if err := s.Refresh(ctx); err != nil {
			return err
		}
		var err error
		resp, err = removeDutsFromDrones(ctx, s, req)
		if err != nil {
			return err
		}
		if err := commitRemoveDuts(ctx, s, resp); err != nil {
			return err
		}
		return nil
	}
	err = retry.Retry(ctx, transientErrorRetries(), f, retry.LogCallback(ctx, "removeDutsFromDronesNoRetry"))
	return resp, err
}

// invCache wraps an InventoryStore and keeps various lookup caches.
type invCache struct {
	store           *gitstore.InventoryStore
	hostnameToID    map[string]string
	droneForDUT     map[string]*inventory.Server
	idToDUT         map[string]*inventory.DeviceUnderTest
	hostnameToDrone map[string]*inventory.Server
}

func newInvCache(ctx context.Context, s *gitstore.InventoryStore) *invCache {
	env := config.Get(ctx).Inventory.Environment
	ic := invCache{
		store:           s,
		hostnameToID:    make(map[string]string),
		droneForDUT:     make(map[string]*inventory.Server),
		idToDUT:         make(map[string]*inventory.DeviceUnderTest),
		hostnameToDrone: make(map[string]*inventory.Server),
	}
	for _, d := range s.Lab.GetDuts() {
		c := d.GetCommon()
		ic.hostnameToID[c.GetHostname()] = c.GetId()
		ic.idToDUT[c.GetId()] = d
	}
	for _, srv := range s.Infrastructure.GetServers() {
		if !isDrone(srv) {
			continue
		}
		if srv.GetEnvironment().String() != env {
			for _, d := range srv.DutUids {
				ic.purgeDUT(d)
			}
			continue
		}
		ic.hostnameToDrone[srv.GetHostname()] = srv
		for _, d := range srv.DutUids {
			ic.droneForDUT[d] = srv
		}
	}
	return &ic
}

// purgeDUT purges the existence of the DUT from the inventory cache.
func (ic *invCache) purgeDUT(dutID string) {
	if d, ok := ic.idToDUT[dutID]; ok {
		delete(ic.hostnameToID, d.GetCommon().GetHostname())
	}
	delete(ic.droneForDUT, dutID)
	delete(ic.idToDUT, dutID)
}

// dutAssigner wraps an InventoryStore and implements assigning DUTs
// to drones.  This struct contains various internal lookup caches.
type dutAssigner struct {
	*invCache
}

func newDUTAssigner(ctx context.Context, s *gitstore.InventoryStore) *dutAssigner {
	return &dutAssigner{
		invCache: newInvCache(ctx, s),
	}
}

func (da *dutAssigner) dutHostnameExists(hostname string) bool {
	_, ok := da.hostnameToID[hostname]
	return ok
}

func (da *dutAssigner) assignDUT(ctx context.Context, a *fleet.AssignDutsToDronesRequest_Item) (*fleet.AssignDutsToDronesResponse_Item, error) {
	ar, err := da.unpackRequest(a)
	if err != nil {
		return nil, err
	}
	if ar.drone == "" {
		ar.drone = pickDroneForDUT(ctx, da.store.Infrastructure)
		logging.Debugf(ctx, "Picked drone %s for DUT %s", ar.drone, a.DutId)
	}
	if _, ok := da.idToDUT[ar.dutID]; !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("DUT %s does not exist", ar.dutID))
	}
	if server, ok := da.droneForDUT[ar.dutID]; ok {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"dut %s is already assigned to drone %s",
			ar.dutID, server.GetHostname(),
		)
	}

	server, ok := da.hostnameToDrone[ar.drone]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("drone %s does not exist", ar.drone))
	}
	server.DutUids = append(server.DutUids, ar.dutID)
	dut := da.idToDUT[ar.dutID]
	dut.RemovalReason = nil

	return &fleet.AssignDutsToDronesResponse_Item{
		DroneHostname: ar.drone,
		DutId:         ar.dutID,
	}, nil
}

// assignRequest is an unpacked fleet.AssignDutsToDronesRequest_Item.
type assignRequest struct {
	dutID string
	drone string
}

func (da *dutAssigner) unpackRequest(r *fleet.AssignDutsToDronesRequest_Item) (assignRequest, error) {
	rr := assignRequest{
		drone: r.DroneHostname,
	}
	if err := da.unpackRequestDUTID(r, &rr); err != nil {
		return rr, err
	}
	return rr, nil
}

func (da *dutAssigner) unpackRequestDUTID(r *fleet.AssignDutsToDronesRequest_Item, rr *assignRequest) error {
	switch {
	case r.DutHostname != "":
		var ok bool
		rr.dutID, ok = da.hostnameToID[r.DutHostname]
		if !ok {
			return status.Errorf(codes.NotFound, "unknown DUT hostname %s", r.DutHostname)
		}
	case r.DutId != "":
		rr.dutID = r.DutId
	default:
		return status.Errorf(codes.InvalidArgument, "must supply one of DUT hostname or ID")
	}
	return nil
}

// pickDroneForDUT returns hostname of a drone to use for the DUT with the given ID.
//
// Returns "" if no drone can be found.
func pickDroneForDUT(ctx context.Context, infra *inventory.Infrastructure) string {
	ds := filterSkylabDronesInEnvironment(ctx, infra.GetServers())
	// ds is a new slice. Sorting it does not modify infra.
	sortDronesByAscendingDUTCount(ds)

	dc := minInt(leastLoadedDronesCount, len(ds))
	ds = ds[:dc]
	if len(ds) == 0 {
		return ""
	}
	return ds[rand.Intn(len(ds))].GetHostname()
}

// commitRemoveDuts commits an in-progress response returned from
// removeDutsFromDrones.
func commitRemoveDuts(ctx context.Context, s *gitstore.InventoryStore, resp *fleet.RemoveDutsFromDronesResponse) error {
	if len(resp.Removed) == 0 {
		return nil
	}
	var err error
	resp.Url, err = s.Commit(ctx, "remove DUTs")
	return err
}

// removeDutsFromDrones implements removing DUTs from drones on an
// InventoryStore.  This is called within a load/commit/retry context.
func removeDutsFromDrones(ctx context.Context, s *gitstore.InventoryStore, req *fleet.RemoveDutsFromDronesRequest) (*fleet.RemoveDutsFromDronesResponse, error) {
	removed := make([]*fleet.RemoveDutsFromDronesResponse_Item, 0, len(req.Removals))
	dr := newDUTRemover(ctx, s)
	for _, r := range req.Removals {
		i, err := dr.removeDUT(ctx, r)
		if err != nil {
			return nil, err
		}
		if i == nil {
			// DUT did not belong to any drone.
			continue
		}
		removed = append(removed, i)
	}
	return &fleet.RemoveDutsFromDronesResponse{
		Removed: removed,
	}, nil
}

// dutRemover wraps an InventoryStore and implements removing DUTs
// from drones.  This struct contains various internal lookup caches.
type dutRemover struct {
	*invCache
}

func newDUTRemover(ctx context.Context, s *gitstore.InventoryStore) *dutRemover {
	return &dutRemover{
		invCache: newInvCache(ctx, s),
	}
}

// removeDUT removes a DUT per a DUT removal request and returns a response.
func (dr *dutRemover) removeDUT(ctx context.Context, r *fleet.RemoveDutsFromDronesRequest_Item) (*fleet.RemoveDutsFromDronesResponse_Item, error) {
	rr, err := dr.unpackRequest(r)
	if err != nil {
		return nil, err
	}
	srv, err := dr.findDroneForRequestDUT(rr)
	if err != nil || srv == nil {
		return nil, err
	}
	srv.DutUids = removeSliceString(srv.DutUids, rr.dutID)
	delete(dr.droneForDUT, rr.dutID)
	d := dr.idToDUT[rr.dutID]
	d.RemovalReason = rr.reason
	return &fleet.RemoveDutsFromDronesResponse_Item{
		DutId:         rr.dutID,
		DroneHostname: srv.GetHostname(),
	}, nil
}

// removeRequest is an unpacked fleet.RemoveDutsFromDronesRequest_Item.
type removeRequest struct {
	dutID  string
	drone  string
	reason *inventory.RemovalReason
}

func (dr *dutRemover) unpackRequest(r *fleet.RemoveDutsFromDronesRequest_Item) (removeRequest, error) {
	rr := removeRequest{
		drone: r.DroneHostname,
	}
	if err := dr.unpackRequestDUTID(r, &rr); err != nil {
		return rr, err
	}
	if err := dr.unpackRequestReason(r, &rr); err != nil {
		return rr, err
	}
	return rr, nil
}

func (dr *dutRemover) unpackRequestDUTID(r *fleet.RemoveDutsFromDronesRequest_Item, rr *removeRequest) error {
	switch {
	case r.DutHostname != "":
		var ok bool
		rr.dutID, ok = dr.hostnameToID[r.DutHostname]
		if !ok {
			return status.Errorf(codes.NotFound, "unknown DUT hostname %s", r.DutHostname)
		}
	case r.DutId != "":
		rr.dutID = r.DutId
	default:
		return status.Errorf(codes.InvalidArgument, "must supply one of DUT hostname or ID")
	}
	return nil
}

func (dr *dutRemover) unpackRequestReason(r *fleet.RemoveDutsFromDronesRequest_Item, rr *removeRequest) error {
	enc := r.GetRemovalReason()
	if len(enc) == 0 {
		return nil
	}
	rr.reason = new(inventory.RemovalReason)
	if err := proto.Unmarshal(enc, rr.reason); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid RemovalReason")
	}
	return nil
}

// findDroneForRequestDUT finds the drone for the DUT in the remove
// request.  If the request listed the wrong drone, return an error.
// Note that this function may return a nil server along with a nil
// error, if the DUT is not found but the request should not consider
// this an error.
func (dr *dutRemover) findDroneForRequestDUT(rr removeRequest) (*inventory.Server, error) {
	srv, ok := dr.droneForDUT[rr.dutID]
	if !ok {
		if rr.drone != "" {
			return nil, status.Errorf(codes.FailedPrecondition, "DUT %s is not assigned to a drone", rr.dutID)
		}
		return nil, nil
	}
	if rr.drone != "" && rr.drone != srv.GetHostname() {
		return nil, status.Errorf(codes.FailedPrecondition, "DUT %s is not on drone %s", rr.dutID, rr.drone)
	}
	return srv, nil
}

// filterSkylabDronesInEnvironment returns drones in the current environment
// from a list of servers
func filterSkylabDronesInEnvironment(ctx context.Context, servers []*inventory.Server) []*inventory.Server {
	env := config.Get(ctx).Inventory.Environment
	ds := make([]*inventory.Server, 0, len(servers))
	for _, s := range servers {
		if s.GetEnvironment().String() != env {
			continue
		}
		for _, r := range s.GetRoles() {
			if r == inventory.Server_ROLE_SKYLAB_DRONE {
				ds = append(ds, s)
				break
			}
		}
	}
	return ds
}

func sortDronesByAscendingDUTCount(ds []*inventory.Server) {
	sort.SliceStable(ds, func(i, j int) bool {
		return len(ds[i].DutUids) < len(ds[j].DutUids)
	})
}

func removeSliceString(sl []string, s string) []string {
	for i, v := range sl {
		if v != s {
			continue
		}
		copy(sl[i:], sl[i+1:])
		return sl[:len(sl)-1]
	}
	return sl
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
