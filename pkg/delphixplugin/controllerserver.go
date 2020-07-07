/*
Package delphixdriver - Delphix CSI Driver
Mainteiner: Daniel Stolf
*/
package delphixdriver

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	deviceID           = "deviceID"
	maxStorageCapacity = tib
)

type accessType int

const (
	mountAccess accessType = iota
	blockAccess
)

type controllerServer struct {
	caps   []*csi.ControllerServiceCapability
	nodeID string
}

// NewControllerServer returns a ControllerServer
func NewControllerServer(nodeID string) *controllerServer {
	return &controllerServer{
		caps: getControllerServiceCapabilities(
			[]csi.ControllerServiceCapability_RPC_Type{
				csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
				csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
				csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
				csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
				csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
			}),
		nodeID: nodeID,
	}
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if err := cs.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("invalid create volume req: %v", req)
		return nil, err
	}

	// Check arguments
	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}
	caps := req.GetVolumeCapabilities()
	if caps == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capabilities missing in request")
	}

	pvcName := req.GetParameters()["csi.storage.k8s.io/pvc/name"]
	namespace := req.GetParameters()["csi.storage.k8s.io/pvc/namespace"]

	// Check for maximum available capacity
	// does this even make sens for Delphix?
	capacity := int64(req.GetCapacityRange().GetRequiredBytes())
	if capacity >= maxStorageCapacity {
		return nil, status.Errorf(codes.OutOfRange, "Requested capacity %d exceeds maximum allowed %d", capacity, maxStorageCapacity)
	}

	// Need to check for already existing volume name
	if exVol, err := getVolumeByName(req.GetName()); err == nil {

		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:      exVol.VolID,
				CapacityBytes: int64(exVol.VolSize),
				VolumeContext: map[string]string{
					"VolID":                exVol.VolID,
					"VolName":              exVol.VolName,
					"VolPath":              exVol.VolPath,
					"VolDelphixReference":  exVol.VolDelphixReference,
					"VolDelphixMountPath":  exVol.VolDelphixMountPath,
					"VolDelphixExportPath": exVol.VolDelphixExportPath,
				},
			},
		}, nil
	}
	// path := getVolumePath(volumeID)

	sourceRef := ""
	ssTimestamp := ""

	if req.GetVolumeContentSource() != nil {

		contentSource := req.GetVolumeContentSource()
		if snapshot := contentSource.GetSnapshot(); snapshot != nil {
			stringSlice := strings.Split(snapshot.GetSnapshotId(), "/")
			sourceRef = stringSlice[0]
			ssTimestamp = stringSlice[2]
		} else if srcVolume := contentSource.GetVolume(); srcVolume != nil {
			glog.Infof("Cloning from Volume ID: \"%v\"", srcVolume.GetVolumeId())
			dv, err := client.FindAppDataContainerByReference(srcVolume.GetVolumeId())
			if err != nil {
				return nil, fmt.Errorf("Can't find Delphix Volume %v: %v", srcVolume.GetVolumeId(), err)
			}
			sourceRef = dv.Reference
			ssTimestamp = "LATEST_SNAPSHOT"
		}
	} else if _, ok := req.Parameters["sourceReference"]; ok {
		sourceRef = req.Parameters["sourceReference"]
		glog.Infof("Cloning from Delphix Object Reference: \"%v\"", sourceRef)
		if _, ok := req.Parameters["sourceSnapshotReference"]; ok {
			ssTimestamp = req.Parameters["sourceSnapshotReference"]
		} else {
			ssTimestamp = "LATEST_SNAPSHOT"
		}
	}

	vol, err := createDriverVolume(req.GetName(), pvcName, namespace, sourceRef, ssTimestamp, capacity /* ephemeral */)
	if err != nil {
		fmt.Printf("ERROR %v \n", err)
		return nil, status.Errorf(codes.Internal, "failed to create volume %v: %v", req.GetName(), err)
	}

	glog.V(4).Infof("created volume %s at path %s", vol.VolID, vol.VolPath)

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      vol.VolID,
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
			VolumeContext: map[string]string{
				"VolID":                  vol.VolID,
				"VolName":                vol.VolName,
				"VolPath":                vol.VolPath,
				"VolDelphixReference":    vol.VolDelphixReference,
				"SourceDelphixReference": vol.SourceDelphixReference,
				"SourceDelphixTimestamp": vol.SourceDelphixTimestamp,
				"VolDelphixMountPath":    vol.VolDelphixMountPath,
				"VolDelphixExportPath":   vol.VolDelphixExportPath,
				"VolDelphixGUID":         vol.VolDelphixGUID,
			},
			ContentSource: req.GetVolumeContentSource(),
		},
	}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	if err := cs.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("invalid delete volume req: %v", req)
		return nil, err
	}

	volID := req.GetVolumeId()
	glog.Infof("Deleting Volume ID: \"%v\"", volID)
	if err := deleteDriverVolume(volID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete volume %v: %v", volID, err)
	}

	glog.Infof("volume %v successfully deleted", volID)

	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.caps,
	}, nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}
	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, req.VolumeId)
	}

	if _, err := getVolumeByID(req.GetVolumeId()); err != nil {
		return nil, status.Error(codes.NotFound, req.GetVolumeId())
	}

	for _, cap := range req.GetVolumeCapabilities() {
		if cap.GetMount() == nil && cap.GetBlock() == nil {
			return nil, status.Error(codes.InvalidArgument, "cannot have both mount and block access type be undefined")
		}

		// A real driver would check the capabilities of the given volume with
		// the set of requested capabilities.
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeContext:      req.GetVolumeContext(),
			VolumeCapabilities: req.GetVolumeCapabilities(),
			Parameters:         req.GetParameters(),
		},
	}, nil
}

func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *controllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *controllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *controllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// getSnapshotPath returns the full path to where the snapshot is stored
func getSnapshotPath(snapshotID string) string {
	return filepath.Join(dataRoot, fmt.Sprintf("%s.tgz", snapshotID))
}

// CreateSnapshot creates a Snapshot request to the Delphix Engine.
func (cs *controllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	if err := cs.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); err != nil {
		glog.V(3).Infof("invalid create snapshot req: %v", req)
		return nil, err
	}

	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}
	// Check arguments
	if len(req.GetSourceVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "SourceVolumeId missing in request")
	}

	// Need to check for already existing snapshot name, and if found check for the
	// requested sourceVolumeId and sourceVolumeId of snapshot that has been created.
	if exSnap, err := getSnapshotByName(req.GetName()); err == nil {
		// Since err is nil, it means the snapshot with the same name already exists need
		// to check if the sourceVolumeId of existing snapshot is the same as in new request.
		if exSnap.VolID == req.GetSourceVolumeId() {
			// same snapshot has been created.
			return &csi.CreateSnapshotResponse{
				Snapshot: &csi.Snapshot{
					SnapshotId:     exSnap.ID,
					SourceVolumeId: exSnap.VolID,
					CreationTime:   &exSnap.CreationTime,
					SizeBytes:      exSnap.SizeBytes,
					ReadyToUse:     exSnap.ReadyToUse,
					// VolumeContext:  exSnap.VolumeContext,
				},
			}, nil
		}
		return nil, status.Errorf(codes.AlreadyExists, "snapshot with the same name: %s but with different SourceVolumeId already exist", req.GetName())
	}

	volumeID := req.GetSourceVolumeId()
	// volumeID := req.GetSourceVolumeId()
	// _, ok := driverVolumes[volumeID]

	// if !ok {
	// 	return nil, status.Error(codes.Internal, "volumeID is not exist")
	// }

	snapshot, err := createDriverSnapshot(volumeID, req.GetName())
	if err != nil {
		return nil, fmt.Errorf("Error creating Snapshot %v , error %v", req.GetName(), err)
	}

	glog.Infof("Create Delphix Volume Snapshot, reference %s", snapshot.SnapshotDelphixReference)

	driverVolumeSnapshots[snapshot.ID] = snapshot

	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     snapshot.ID,
			SourceVolumeId: snapshot.VolID,
			CreationTime:   &snapshot.CreationTime,
			SizeBytes:      snapshot.SizeBytes,
			ReadyToUse:     snapshot.ReadyToUse,
		},
	}, nil
}

func (cs *controllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	// Check arguments
	if len(req.GetSnapshotId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Snapshot ID missing in request")
	}

	if err := cs.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); err != nil {
		glog.V(3).Infof("invalid delete snapshot req: %v", req)
		return nil, err
	}
	snapshotID := req.GetSnapshotId()
	glog.V(4).Infof("deleting snapshot %s", snapshotID)
	path := getSnapshotPath(snapshotID)
	os.RemoveAll(path)
	delete(driverVolumeSnapshots, snapshotID)
	return &csi.DeleteSnapshotResponse{}, nil
}

func (cs *controllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	if err := cs.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS); err != nil {
		glog.V(3).Infof("invalid list snapshot req: %v", req)
		return nil, err
	}

	// case 1: SnapshotId is not empty, return snapshots that match the snapshot id.
	if len(req.GetSnapshotId()) != 0 {
		snapshotID := req.SnapshotId
		if snapshot, ok := driverVolumeSnapshots[snapshotID]; ok {
			return convertSnapshot(snapshot), nil
		}
	}

	// case 2: SourceVolumeId is not empty, return snapshots that match the source volume id.
	if len(req.GetSourceVolumeId()) != 0 {
		for _, snapshot := range driverVolumeSnapshots {
			if snapshot.VolID == req.SourceVolumeId {
				return convertSnapshot(snapshot), nil
			}
		}
	}

	var snapshots []csi.Snapshot
	// case 3: no parameter is set, so we return all the snapshots.
	sortedKeys := make([]string, 0)
	for k := range driverVolumeSnapshots {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		snap := driverVolumeSnapshots[key]
		snapshot := csi.Snapshot{
			SnapshotId:     snap.ID,
			SourceVolumeId: snap.VolID,
			CreationTime:   &snap.CreationTime,
			SizeBytes:      snap.SizeBytes,
			ReadyToUse:     snap.ReadyToUse,
		}
		snapshots = append(snapshots, snapshot)
	}

	var (
		ulenSnapshots = int32(len(snapshots))
		maxEntries    = req.MaxEntries
		startingToken int32
	)

	if v := req.StartingToken; v != "" {
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, status.Errorf(
				codes.Aborted,
				"startingToken=%d !< int32=%d",
				startingToken, math.MaxUint32)
		}
		startingToken = int32(i)
	}

	if startingToken > ulenSnapshots {
		return nil, status.Errorf(
			codes.Aborted,
			"startingToken=%d > len(snapshots)=%d",
			startingToken, ulenSnapshots)
	}

	// Discern the number of remaining entries.
	rem := ulenSnapshots - startingToken

	// If maxEntries is 0 or greater than the number of remaining entries then
	// set maxEntries to the number of remaining entries.
	if maxEntries == 0 || maxEntries > rem {
		maxEntries = rem
	}

	var (
		i       int
		j       = startingToken
		entries = make(
			[]*csi.ListSnapshotsResponse_Entry,
			maxEntries)
	)

	for i = 0; i < len(entries); i++ {
		entries[i] = &csi.ListSnapshotsResponse_Entry{
			Snapshot: &snapshots[j],
		}
		j++
	}

	var nextToken string
	if j < ulenSnapshots {
		nextToken = fmt.Sprintf("%d", j)
	}

	return &csi.ListSnapshotsResponse{
		Entries:   entries,
		NextToken: nextToken,
	}, nil
}

func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {

	volID := req.GetVolumeId()
	if len(volID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	capRange := req.GetCapacityRange()
	if capRange == nil {
		return nil, status.Error(codes.InvalidArgument, "Capacity range not provided")
	}

	capacity := int64(capRange.GetRequiredBytes())
	if capacity >= maxStorageCapacity {
		return nil, status.Errorf(codes.OutOfRange, "Requested capacity %d exceeds maximum allowed %d", capacity, maxStorageCapacity)
	}

	exVol, err := getVolumeByID(volID)
	if err != nil {
		// Assume not found error
		return nil, status.Errorf(codes.NotFound, "Could not get volume %s: %v", volID, err)
	}

	if exVol.VolSize < capacity {
		exVol.VolSize = capacity
		if err := updateDriverVolume(volID, exVol); err != nil {
			return nil, status.Errorf(codes.Internal, "Could not update volume %s: %v", volID, err)
		}
	}

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         exVol.VolSize,
		NodeExpansionRequired: true,
	}, nil
}

func convertSnapshot(snap driverSnapshot) *csi.ListSnapshotsResponse {
	entries := []*csi.ListSnapshotsResponse_Entry{
		{
			Snapshot: &csi.Snapshot{
				SnapshotId:     snap.ID,
				SourceVolumeId: snap.VolID,
				CreationTime:   &snap.CreationTime,
				SizeBytes:      snap.SizeBytes,
				ReadyToUse:     snap.ReadyToUse,
			},
		},
	}

	rsp := &csi.ListSnapshotsResponse{
		Entries: entries,
	}

	return rsp
}

func (cs *controllerServer) validateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, cap := range cs.caps {
		if c == cap.GetRpc().GetType() {
			return nil
		}
	}
	return status.Errorf(codes.InvalidArgument, "unsupported capability %s", c)
}

func getControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) []*csi.ControllerServiceCapability {
	var csc []*csi.ControllerServiceCapability

	for _, cap := range cl {
		glog.Infof("Enabling controller service capability: %v", cap.String())
		csc = append(csc, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		})
	}

	return csc
}
