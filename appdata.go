package delphix

import (
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// FindAppDataContainerByReference transforms an interface into an AppDataContainerStruct
func (c *Client) FindAppDataContainerByReference(r string) (AppDataContainerStruct, error) {
	var appDataObj AppDataContainerStruct
	var err error
	dbObj, err := c.FindDatabaseByReference(r)
	if err != nil {
		return appDataObj, err
	}
	t := dbObj.(map[string]interface{})["type"].(string)

	if t == "AppDataContainer" {
		err = mapstructure.Decode(dbObj, &appDataObj)
	} else {
		err = fmt.Errorf("ERROR: Type was %s, expected AppDataContainer", t)
	}
	return appDataObj, err
}

// FindAppDataContainerByName transforms an interface into an AppDataContainerStruct
func (c *Client) FindAppDataContainerByName(n string) (AppDataContainerStruct, error) {
	var appDataObj AppDataContainerStruct
	var err error
	dbObj, err := c.FindDatabaseByName(n)
	if err != nil || dbObj == interface{} {
		return appDataObj, err
	}

	t := dbObj.(map[string]interface{})["type"].(string)

	if t == "AppDataContainer" {
		err = mapstructure.Decode(dbObj, &appDataObj)
	} else {
		err = fmt.Errorf("ERROR: Type was %s, expected AppDataContainer", t)
	}
	return appDataObj, err
}

// FindAppDataRepoByEnvironmentRef returns the reference of the Repository by
// environment reference (e)
func (c *Client) FindAppDataRepoByEnvironmentRef(r string) (string, error) {
	var repoRef string
	result, err := c.FindObject("repository", "environment", r, 0)
	if err != nil {
		log.Fatalf("Wamp, Wamp: %s\n", err)
	} else if result == nil {
		log.Fatalf("[ERROR] Repo for Environment Reference %s not found\n", r)
		err = fmt.Errorf("ERROR: Repository not found for Environment %s", r)
	} else {
		repoRef = result.(map[string]interface{})["reference"].(string)
	}
	return repoRef, err //return the reference
}

// FindLastAppDataSnapshotByTimeflow returns the snapshot (AppDataSnapshotStruct) of the
// last Snapshot in a Timeflow
func (c *Client) FindLastAppDataSnapshotByTimeflow(r string) (AppDataSnapshotStruct, error) {
	var appDatasnapObj AppDataSnapshotStruct
	var err error
	snapObj, err := c.FindLastSnapshotByTimeflow(r)
	if err != nil {
		return appDatasnapObj, err
	}
	t := snapObj.(map[string]interface{})["type"].(string)

	if t == "AppDataSnapshot" {
		err = mapstructure.Decode(snapObj, &appDatasnapObj)
	} else {
		err = fmt.Errorf("ERROR: Type was %s, expected AppDataSnapshot", t)
	}
	return appDatasnapObj, err
}

// FindAppDataSnapshotByReference returns the snapshot (interface{}) searching by Reference
func (c *Client) FindAppDataSnapshotByReference(r string) (AppDataSnapshotStruct, error) {
	var appDatasnapObj AppDataSnapshotStruct
	var err error
	snapObj, err := c.FindSnapshotByReference(r)
	if err != nil {
		return appDatasnapObj, err
	}
	t := snapObj.(map[string]interface{})["type"].(string)

	if t == "AppDataSnapshot" {
		err = mapstructure.Decode(snapObj, &appDatasnapObj)
	} else {
		err = fmt.Errorf("ERROR: Type was %s, expected AppDataSnapshot", t)
	}
	return appDatasnapObj, err
}

// ProvisionCloneAppData provisions an empty vFile Folder
func (c *Client) ProvisionCloneAppData(v *AppDataProvisionParametersStruct) (interface{}, error) {
	reference, err := c.executePostJobAndReturnObjectReference("/database/provision", v)
	return reference, err
}

// ProvisionAppData provisions an AppData Container (clone or empty)
// name - name of the new volume
// envRef - environment reference
// groupRef - group reference the AppData will belong to. Defaults to GROUP-1
// sourceRef - dsource reference. will provision new empty AppData if true
// timestamp - delphix timestamp or LATEST_SNAPSHOT. Defaults to LATEST_SNAPSHOT
func (c *Client) ProvisionAppData(name string, path string, envRef string, groupRef string, sourceRef string, timestamp string) (interface{}, error) {
	fbool := new(bool)
	*fbool = false

	containerObj := AppDataContainerStruct{
		Group: groupRef,
		Name:  name,
		Type:  "AppDataContainer",
	}
	// sourceObj, err := c.FindAppDataContainerByReference(sourceRef)
	repoRef, _ := c.FindAppDataRepoByEnvironmentRef(envRef)

	sourceStructObj := &AppDataVirtualSourceStruct{
		Type:                            "AppDataVirtualSource",
		AllowAutoVDBRestartOnHostReboot: fbool,
		Name:                            name,
		AdditionalMountPoints:           []*AppDataAdditionalMountPointStruct{},
		Operations: &VirtualSourceOperationsStruct{
			Type: "VirtualSourceOperations",
		},
		Parameters: &JSONStruct{},
	}
	userRef, _ := c.FindEnvironmentPrimaryUser(envRef)
	sourceConfigObj := &AppDataSourceConfigStruct{
		Name:            path,
		Repository:      repoRef,
		EnvironmentUser: userRef,
		Path:            path,
		Parameters:      &JSONStruct{},
		Type:            "AppDataDirectSourceConfig",
	}

	timeflowParameter := &AppDataTimeflowPointParametersStruct{}
	if timestamp == "" || timestamp == "LATEST_SNAPSHOT" {
		timeflowParameter = &AppDataTimeflowPointParametersStruct{
			Type:      "TimeflowPointSemantic",
			Container: sourceRef,
			Location:  "LATEST_SNAPSHOT",
		}
	} else {
		timeflowParameter = &AppDataTimeflowPointParametersStruct{
			Type:      "TimeflowPointSemantic",
			Container: sourceRef,
			Timestamp: timestamp,
		}
	}

	provisionParametersObj := AppDataProvisionParametersFactory(&containerObj, fbool, "", sourceStructObj, sourceConfigObj, &timeflowParameter)
	got, err := c.ProvisionCloneAppData(&provisionParametersObj)
	if err != nil || got == nil {
		log.Println("[ERROR] Client.ProvisionCloneAppData() error = ", err)
	}
	return got, err

}

// DeleteAppData permanently deletes an AppData VDB
func (c *Client) DeleteAppData(v string) error {

	o := DeleteParametersStruct{
		Type: "DeleteParameters",
	}
	url := fmt.Sprintf("/database/%s/delete", strings.ToUpper(v))
	err := c.executePostJobAndReturnErrOnly(url, o)
	return err
}

// SyncAppData performs a snapsync on an AppData VDB
func (c *Client) SyncAppData(r string) (AppDataSnapshotStruct, error) {
	var snapObj AppDataSnapshotStruct
	// check if AppData Reference is valid
	appDataObj, _ := c.FindAppDataContainerByReference(r)
	if (appDataObj == AppDataContainerStruct{}) {
		log.Println("[ERROR] AppData", r, "not found")
		err := fmt.Errorf("AppData %s not found", r)
		return snapObj, err
	}

	url := fmt.Sprintf("/database/%s/sync", strings.ToUpper(r))
	appDataSyncParameters := CreateAppDataSyncParameters(nil)
	err := c.executePostJobAndReturnErrOnly(url, appDataSyncParameters)
	if err != nil {
		log.Println("[ERROR] Something went wrong executing Snapshot: ", err.Error())
		return snapObj, err
	}
	appDataObj, err = c.FindAppDataContainerByReference(r)
	snapObj, err = c.FindLastAppDataSnapshotByTimeflow(appDataObj.CurrentTimeflow)

	return snapObj, err
}

// RollbackAppData rollbacks an AppDataContainer r to AppDataSnapshot s
func (c *Client) RollbackAppData(d string, s string) (AppDataSnapshotStruct, error) {
	url := fmt.Sprintf("/database/%s/rollback", strings.ToUpper(d))
	var response AppDataSnapshotStruct
	appDataObj, _ := c.FindAppDataContainerByReference(d)
	snapObj, _ := c.FindAppDataSnapshotByReference(s)
	if (appDataObj == AppDataContainerStruct{} || snapObj == AppDataSnapshotStruct{}) {
		err := fmt.Errorf("No Snapshot %s found for AppData %s", s, d)
		return response, err
	}
	timestamp := strings.TrimPrefix(snapObj.Name, "@") + "z"
	rollbackParameters := CreateAppDataRollbackParameters(timestamp, appDataObj.CurrentTimeflow)
	err := c.executePostJobAndReturnErrOnly(url, rollbackParameters)
	if err != nil {
		log.Println("[ERROR] Something went wrong with AppData rollback: ", err.Error())
		return response, err
	}
	// After rollback, the Timeflow has changed, so we need to update the AppData Object
	// TO-DO - this is NOT working... looks like Delphix takes a few seconds to update the Timeflow, event after job is done
	appDataObj, err = c.FindAppDataContainerByReference(d)
	response, err = c.FindLastAppDataSnapshotByTimeflow(appDataObj.CurrentTimeflow)
	return response, err
}
