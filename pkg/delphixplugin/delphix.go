/*
Package delphixdriver - Delphix CSI Driver
Mainteiner: Daniel Stolf
*/
package delphixdriver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	delphix "github.com/DSTOLF/delphix-go-sdk"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
)

type driver struct {
	name            string
	nodeID          string
	client          delphix.Client
	endpoint        string
	environmentName string
	groupName       string
	repoName        string
	sourceRef       string
	mountPath       string
	version         string

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer
}

type driverVolume struct {
	VolName                string     `json:"volName"`
	VolID                  string     `json:"volID"`
	VolSize                int64      `json:"volSize"`
	VolPath                string     `json:"volPath"`
	VolDelphixReference    string     `json:"volDelphixReference"`
	SourceDelphixReference string     `json:"sourceDelphixReference"`
	SourceDelphixTimestamp string     `json:"sourceDelphixTimestamp"`
	VolDelphixMountPath    string     `json:"volDelphixMountPath"`
	VolDelphixExportPath   string     `json:"VolDelphixGUID"`
	VolDelphixGUID         string     `json:"VolDelphixExportPath"`
	VolAccessType          accessType `json:"volAccessType"`
}

type driverSnapshot struct {
	Name                     string              `json:"name"`
	ID                       string              `json:"id"`
	VolID                    string              `json:"volID"`
	CreationTime             timestamp.Timestamp `json:"creationTime"`
	SnapshotDelphixReference string              `json:"snapshotDelphixReference"`
	SizeBytes                int64               `json:"sizeBytes"`
	VolDelphixReference      string              `json:"volDelphixReference"`
	ReadyToUse               bool                `json:"readyToUse"`
}

// MountPointRequest is a struct for the requests this server will receive
// this is just for the ExportPath workaround
type MountPointRequest struct {
	MountPath string `json:"mountpath"`
}

var (
	vendorVersion         = "dev"
	client                = &delphix.Client{}
	envRef                string
	groupRef              string
	mountPath             string
	sourceRef             string
	repoName              string
	driverVolumes         map[string]driverVolume
	driverVolumeSnapshots map[string]driverSnapshot
)

const (
	// Directory where data for volumes and snapshots are persisted.
	// This can be ephemeral within the container or persisted if
	// backed by a Pod volume.
	dataRoot = "/csi-data-dir"
)

func init() {
	driverVolumes = map[string]driverVolume{}
	driverVolumeSnapshots = map[string]driverSnapshot{}
}

// NewDriver returns a driver object
func NewDriver(driverName, nodeID, endpoint string, delphixClient delphix.Client, environmentName, groupName, repName, srcRef, mPath, version string) (*driver, error) {
	env, err := delphixClient.FindEnvironmentByName(environmentName)
	group, err := delphixClient.FindGroupByName(groupName)
	if err != nil || env == nil {
		return nil, fmt.Errorf("Environment doesn't exist in Delphix: %v", err)
	}

	envRef = env.(map[string]interface{})["reference"].(string)
	groupRef = group.(map[string]interface{})["reference"].(string)
	repoName = repName
	mountPath = mPath
	sourceRef = srcRef
	client = &delphixClient
	if driverName == "" {
		return nil, fmt.Errorf("No driver name provided")
	}

	if nodeID == "" {
		return nil, fmt.Errorf("No node id provided")
	}

	if mountPath == "" {
		return nil, fmt.Errorf("No Delphix Mount Path provided")
	}

	if repoName == "" {
		return nil, fmt.Errorf("No Delphix Repository Name provided")
	}

	if repoName == "" {
		return nil, fmt.Errorf("No Empty Source Ref provided")
	}

	if version != "" {
		vendorVersion = version
	}

	if err := os.MkdirAll(dataRoot, 0750); err != nil {
		return nil, fmt.Errorf("failed to create dataRoot: %v", err)
	}

	glog.Infof("Driver: %v ", driverName)
	glog.Infof("Version: %s", vendorVersion)

	return &driver{
		name:            driverName,
		version:         vendorVersion,
		nodeID:          nodeID,
		client:          delphixClient,
		environmentName: environmentName,
		mountPath:       mountPath,
		endpoint:        endpoint,
	}, nil
}

func (de *driver) Run() {
	// Create GRPC servers
	de.ids = NewIdentityServer(de.name, de.version)
	de.ns = NewNodeServer(de.nodeID, de.client)
	de.cs = NewControllerServer(de.nodeID)

	s := NewNonBlockingGRPCServer()
	s.Start(de.endpoint, de.ids, de.cs, de.ns)
	s.Wait()
}

func getVolumeByID(volumeID string) (driverVolume, error) {
	if driverVol, ok := driverVolumes[volumeID]; ok {
		return driverVol, nil
	}
	return driverVolume{}, fmt.Errorf("volume id %s does not exit in the volumes list", volumeID)
}

func getVolumeByName(volName string) (driverVolume, error) {
	for _, driverVol := range driverVolumes {
		if driverVol.VolName == volName {
			return driverVol, nil
		}
	}
	return driverVolume{}, fmt.Errorf("volume name %s does not exit in the volumes list", volName)
}

func getSnapshotByName(name string) (driverSnapshot, error) {
	for _, snapshot := range driverVolumeSnapshots {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return driverSnapshot{}, fmt.Errorf("snapshot name %s does not exit in the snapshots list", name)
}

// getVolumePath returs the canonical path for Delphix volume
func getVolumePath(volID string) string {
	return filepath.Join(dataRoot, volID)
}

// createVolume create the directory for the Delphix volume.
// It returns the volume path or err if one occurs.

func createDriverVolume(name, pvcName, namespace, srcRef, timestamp string, cap int64) (*driverVolume, error) {
	var src string
	fbool := new(bool)
	*fbool = false
	err := client.Login()

	if err != nil {
		return nil, fmt.Errorf("Failed Logging to Delphix: %v", err)
	}

	// We need to be sure this volume wasn't previously created
	// There's no guarantee Kubernetes won't call this function more than once per request
	// This function needs to be idempotent, so we can't return error if the volume already exists
	if srcRef == "" {
		src = sourceRef
	} else {
		src = srcRef
	}

	delphixObj, err := client.FindCSIAppDataByName(name)
	if err != nil {
		return nil, fmt.Errorf("Failed querying Delphix for Volume: %v", err)
	}
	if (delphixObj == delphix.AppDataContainerStruct{}) {
		ref, err := client.ProvisionCSIAppData(name, envRef, repoName, groupRef, src, timestamp, pvcName, name, namespace, mountPath)
		if (err != nil || ref == delphix.AppDataContainerStruct{}) {
			glog.V(4).Infof("Error Provisioning Volume %v: %v", ref.(string), err)
			//If there's an error, we need to make really sure the volume really isn't there
			_ = client.DeleteAppData(ref.(string))
			return nil, fmt.Errorf("Error creating volume: %v", err)
		}
		delphixObj, err = client.FindCSIAppDataByReference(ref.(string))
	}

	glog.V(4).Infof("Delphix Volume Created: %v", delphixObj.Reference)

	delphixReference := delphixObj.Reference
	volID := delphixReference
	path := getVolumePath(volID)

	delphixGUID := delphixObj.GUID
	appDataSource, err := client.FindSourceByContainerReference(delphixObj.Reference)
	appDataSourceConfigReference := appDataSource.(map[string]interface{})["config"].(string)
	appDataSourceConfig, err := client.FindAppDataCSISourceConfigByReference(appDataSourceConfigReference)
	exportPath := appDataSourceConfig.Parameters.ExportPath

	err = os.MkdirAll(path, 0777)
	if err != nil {
		return nil, err
	}

	delphixVol := driverVolume{
		VolID:                  volID,
		VolName:                name,
		VolSize:                cap,
		VolPath:                path,
		VolDelphixReference:    delphixReference,
		SourceDelphixReference: sourceRef,
		SourceDelphixTimestamp: timestamp,
		VolDelphixMountPath:    mountPath,
		VolDelphixExportPath:   exportPath,
		VolDelphixGUID:         delphixGUID,
	}
	driverVolumes[volID] = delphixVol
	return &delphixVol, nil
}

// createDriversnapshot
func createDriverSnapshot(volID, name string) (driverSnapshot, error) {

	dbObj, err := client.FindAppDataContainerByReference(volID)
	ssObj := &driverSnapshot{}
	// check if volume exists in Delphix
	if err != nil {
		return *ssObj, fmt.Errorf("Unable to find Delphix volume %v, error: %v", volID, err)
	}

	// create a snapshot in Delphix
	ssRef, err := client.SyncAppData(dbObj.Reference)

	// Making a jugdment call here, might need to revisit this decision later
	// Snapshots don't have Context, like a Volume, where we can add all information we need \
	// to interact with the engine later
	// So right now, source Reference, snapshot reference and name (which is a timestamp). Is this enough??
	// snapshotID := uuid.NewUUID().String()

	snapshotID := dbObj.Reference + "/" + ssRef.Reference + "/" + strings.TrimPrefix(ssRef.Name, "@") + "z"

	if err != nil {
		return *ssObj, fmt.Errorf("Can't create snapshot Delphix volume %v, error: %v", volID, err)
	}

	// creationTime := ptypes.TimestampNow()
	creationTime, err := time.Parse(ssRef.LatestChangePoint.Timestamp, ssRef.LatestChangePoint.Timestamp)
	creationTimeStamp, err := ptypes.TimestampProto(creationTime)

	if err != nil {
		return *ssObj, fmt.Errorf("Can't convert timestamp %v, error: %v", ssRef.LatestChangePoint.Timestamp, err)
	}

	tbool := true
	delphixSnapshot := driverSnapshot{
		Name:                     name,
		ID:                       snapshotID,
		VolID:                    volID,
		CreationTime:             *creationTimeStamp,
		SnapshotDelphixReference: ssRef.Reference,
		VolDelphixReference:      dbObj.Reference,
		ReadyToUse:               tbool,
	}
	return delphixSnapshot, nil
}

// updateVolume updates the existing Delphix volume.
func updateDriverVolume(volID string, volume driverVolume) error {
	glog.V(4).Infof("updating Delphix volume: %s", volID)

	if _, err := getVolumeByID(volID); err != nil {
		return err
	}

	driverVolumes[volID] = volume
	return nil
}

// deleteVolume deletes the directory for the Delphix volume.
func deleteDriverVolume(volDelphixReference string) error {
	glog.Infof("Deleting Delphix volume: %s", volDelphixReference)

	_, err := client.FindAppDataContainerByReference(volDelphixReference)

	if err == nil {
		err = client.DeleteAppData(volDelphixReference)
		if err != nil {
			glog.Infof("Error deleting Delphix volume: %s, %s", volDelphixReference, err)
		}
	} else {
		glog.Infof("Delphix Volume Reference doesn't exist: %s", volDelphixReference)
	}

	path := getVolumePath(volDelphixReference)
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	delete(driverVolumes, volDelphixReference)
	return nil
}

// driverIsEmpty is a simple check to determine if the specified target path directory
// is empty or not.
func driverIsEmpty(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		return true, fmt.Errorf("unable to open target path volume, error: %v", err)
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
