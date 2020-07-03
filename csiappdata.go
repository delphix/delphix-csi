package delphix

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"
)

// FindCSIRepoByNameAndEnvironmentRef returns the reference of the Repository by
// environment reference (e)
func (c *Client) FindCSIRepoByNameAndEnvironmentRef(name, envRef string) (string, error) {
	var err error
	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		Get(c.url + "/repository?environment=" + envRef) //grab all the repos for the environment
	if err != nil {
		return "", err
	}

	if http.StatusOK != resp.StatusCode() { //check to make sure our query was good
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		log.Printf(c.url + "/repository?environment=" + envRef)
		log.Printf(errorMessage)
		if err != nil {
			return "", err
		}
	}

	repos := resp.Body()

	var dat map[string]interface{}

	if err = json.Unmarshal(repos, &dat); err != nil { //convert the json to go objects
		return "", err
	}

	results := dat["result"].([]interface{}) //grab the query results
	for _, result := range results {         //loop through the repos
		n := result.(map[string]interface{})["name"].(string) //grab the repo name
		//grab the repo UUID
		if n == name { //if the name matches our specified repo
			reference := result.(map[string]interface{})["reference"].(string)
			return reference, nil //return the reference
		}
	}
	log.Printf("Unable to find repository of type app data on %s", envRef)
	return "", nil
}

// FindCSIAppDataByReference transforms an interface into an AppDataContainerStruct
func (c *Client) FindCSIAppDataByReference(r string) (AppDataContainerStruct, error) {
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

// FindCSIAppDataByName transforms an interface into an AppDataContainerStruct
func (c *Client) FindCSIAppDataByName(n string) (AppDataContainerStruct, error) {
	var appDataObj AppDataContainerStruct
	var err error
	dbObj, err := c.FindDatabaseByName(n)
	if err != nil {
		return appDataObj, err
	}
	if dbObj == nil {
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

// ProvisionCSIAppData provisions an AppData Container (clone or empty)
// name - name of the new volume
// envRef - environment reference
// groupRef - group reference the AppData will belong to. Defaults to GROUP-1
// sourceRef - dsource reference. will provision new empty AppData if true
// timestamp - delphix timestamp or LATEST_SNAPSHOT. Defaults to LATEST_SNAPSHOT
func (c *Client) ProvisionCSIAppData(name string, envRef string, repoName string, groupRef string, sourceRef string, timestamp string, pvcName string, pvName string, namespace, mountPath string) (interface{}, error) {
	fbool := new(bool)
	*fbool = false

	containerObj := AppDataContainerStruct{
		Group: groupRef,
		Name:  name,
		Type:  "AppDataContainer",
	}
	sourceConfigParameter := &CSISourceConfigParameterStruct{
		MountLocation:         mountPath,
		PersistentVolumeClaim: pvcName,
		PersistentVolume:      pvName,
		Namespace:             namespace,
	}
	// sourceObj, err := c.FindAppDataContainerByReference(sourceRef)
	repoRef, _ := c.FindCSIRepoByNameAndEnvironmentRef(repoName, envRef)
	sourceStructObj := &CSIAppDataVirtualSourceStruct{
		Type:                            "AppDataVirtualSource",
		AllowAutoVDBRestartOnHostReboot: fbool,
		Name:                            name,
		AdditionalMountPoints:           []*AppDataAdditionalMountPointStruct{},
		Operations: &VirtualSourceOperationsStruct{
			Type: "VirtualSourceOperations",
		},
		Parameters: sourceConfigParameter,
	}
	userRef, _ := c.FindEnvironmentPrimaryUser(envRef)
	sourceConfigObj := &CSISourceConfigStruct{
		Name:            name,
		Repository:      repoRef,
		EnvironmentUser: userRef,
		Parameters:      sourceConfigParameter,
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

	provisionParametersObj := CSIProvisionParametersFactory(&containerObj, fbool, "", sourceStructObj, sourceConfigObj, &timeflowParameter)

	reference, err := c.executePostJobAndReturnObjectReference("/database/provision", provisionParametersObj)

	//
	if err != nil {
		glog.Fatalf("[ERROR] Client.ProvisionCSIAppData() error = %s", err)
		return reference, err
	}

	if reference == nil {
		err = fmt.Errorf("[ERROR] Client.ProvisionCSIAppData() - nil reference")
		glog.V(2).Infof("[ERROR] Client.ProvisionCSIAppData() - nil reference")
		return reference, err
	}
	return reference.(string), err
	// query source by database ref -> get config
	// query source config
}

// FindAppDataCSISourceConfigByReference returns the SourceConfig object (interface) by reference (r) SourceConfig
func (c *Client) FindAppDataCSISourceConfigByReference(r string) (CSISourceConfigStruct, error) {
	var csiSourceconfigObj CSISourceConfigStruct
	// var csiSourceConfigParameterObj CSISourceConfigParameterStruct
	result, err := c.FindSourceConfigByReference(r)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("SourceConfig Reference %s not found\n", r)
		return csiSourceconfigObj, err
	}

	parameters := result.(map[string]interface{})["parameters"]

	err = mapstructure.Decode(result, &csiSourceconfigObj)

	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}

	exportPath := parameters.(map[string]interface{})["export_path"]
	pvc := parameters.(map[string]interface{})["persistent_volume_claim"]
	pv := parameters.(map[string]interface{})["persistent_volume"]
	if parameters == nil || exportPath == nil || pvc == nil || pv == nil {
		err = fmt.Errorf("[ERROR] Object returned with empty parameters %v ", parameters)
		return CSISourceConfigStruct{}, err
	}

	// hideous workaround because I couldn't get mapstructure or even json.Marshal/Unmarshal to work with the Parameters field
	csiSourceconfigObj.Parameters.ExportPath = exportPath.(string)
	csiSourceconfigObj.Parameters.PersistentVolumeClaim = pvc.(string)
	csiSourceconfigObj.Parameters.PersistentVolume = pv.(string)

	// csiSourceconfigObj.Parameters = csiSourceConfigParameterObj
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}

	return csiSourceconfigObj, err //return the reference
}
