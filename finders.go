package delphix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/glog"
)

// FindObject returns the database object (interface) of the named (n) database
func (c *Client) FindObject(objectType, fieldName, fieldValue string, pageSize int) (interface{}, error) {
	var err error
	var url string
	if fieldName == "reference" { //grab database by reference
		url = c.url + "/" + objectType + "/" + fieldValue
	} else { //query databases for field=value
		url = c.url + "/" + objectType + "?" + fieldName + "=" + fieldValue
		if pageSize > 0 {
			url = url + "&pageSize=" + strconv.Itoa(pageSize)
		}
	}

	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}

	if http.StatusOK != resp.StatusCode() { //check to make sure our query was good
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return nil, err
		}
	}

	data := resp.Body()
	var response map[string]interface{}
	if err = json.Unmarshal(data, &response); err != nil { //convert the json to go objects
		return nil, err
	}

	responseType := response["type"]

	if responseType == "ListResult" {
		results := response["result"].([]interface{}) //grab the query results

		for _, result := range results { //loop through the databases
			r := result.(map[string]interface{})[fieldName] //grab the db name
			// reference := result.(map[string]interface{})["reference"] //grab the db UUID
			if fieldValue == r { //if the name matches our specified vdb
				return result.(map[string]interface{}), nil //return the reference
			}
		}
	} else {
		if err = json.Unmarshal(data, &response); err != nil { //convert the json to go objects
			glog.Fatalf("Error unmarshaling JSON: %s", err)
			return nil, err
		}
		database := response["result"].(interface{})  //grab the query results
		return database.(map[string]interface{}), nil //return the reference
	}
	return nil, nil
}

// FindGroupByReference returns the group object (interface) of the reference (r) group
func (c *Client) FindGroupByReference(r string) (interface{}, error) {
	result, err := c.FindObject("group", "reference", r, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Group Reference %s not found\n", r)
	}
	return result, err //return the reference
}

// FindGroupByName returns the group object (interface) of the name (n) group
func (c *Client) FindGroupByName(n string) (interface{}, error) {
	result, err := c.FindObject("group", "name", n, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		err = fmt.Errorf("ERROR: Group Name %s not found", n)
		glog.V(2).Infof("[ERROR] Group Name %s not found\n", n)
	}
	return result, err //return the reference
}

// FindDatabaseByReference returns the database object (interface) of the reference (r) Database
func (c *Client) FindDatabaseByReference(r string) (interface{}, error) {
	result, err := c.FindObject("database", "reference", r, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Database Reference %s not found\n", r)
	}
	return result, err //return the reference
}

// FindDatabaseByName returns the database object (interface) of the named (n) database
func (c *Client) FindDatabaseByName(n string) (interface{}, error) {
	result, err := c.FindObject("database", "name", n, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Database Name %s not found\n", n)
	}
	return result, err //return the reference
}

// FindSnapshotByReference returns the timestamp (string) of the Snapshot
func (c *Client) FindSnapshotByReference(r string) (interface{}, error) {
	result, err := c.FindObject("snapshot", "reference", r, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Snapshot Reference %s not found\n", r)
	}
	return result, err //return the reference
}

// FindLastSnapshotByTimeflow returns the snapshot (string) of the
// last Snapshot in a Timeflow t
func (c *Client) FindLastSnapshotByTimeflow(t string) (interface{}, error) {
	result, err := c.FindObject("snapshot", "timeflow", t, 1)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Snapshot For Timeflow %s not found\n", t)
	}
	return result, err //return the reference
}

// FindSourceByContainerReference returns the source object (interface) by reference (r) Source (external database)
func (c *Client) FindSourceByContainerReference(r string) (interface{}, error) {
	// this is the only function that can't use the generic find object function
	// the query parameter accepted by the API is "database", but the actual field is named "container" (because... why wouldn't it be, right?)
	// since the response type is ListResult (again: why wouldn't it, right?)
	// we need to iterate over each field, searching for one named "database", which doesn't exit
	// since this is such a critical path for anyone using this for object created by plugins,
	// it was decided to create aaaall these small specific functions for the sake of consistency
	var err error

	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		Get(c.url + "/source?database=" + r) //query source by reference
	if err != nil {
		return nil, err
	}

	if http.StatusOK != resp.StatusCode() { //check to make sure our query was good
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		glog.Fatalf(errorMessage)
		if err != nil {
			return nil, err
		}
	}

	dat := resp.Body()

	var response map[string]interface{}
	if err = json.Unmarshal(dat, &response); err != nil { //convert the json to go objects
		return nil, err
	}
	results := response["result"].([]interface{}) //grab the query results
	for _, database := range results {            //loop through the databases
		reference := database.(map[string]interface{})["container"] //grab the db UUID
		if reference == r {                                         //if the name matches our specified vdb
			return database.(map[string]interface{}), nil //return the reference
		}
	}

	errorMessage := "Unable to find Source with a Container Reference of " + r + " in " + string(resp.Body())
	err = fmt.Errorf(errorMessage)
	return nil, err
}

// FindSourceConfigByReference returns the SourceConfig object (interface) by reference (r) SourceConfig
func (c *Client) FindSourceConfigByReference(r string) (interface{}, error) {
	result, err := c.FindObject("sourceconfig", "reference", r, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("SourceConfig Reference %s not found\n", r)
	}
	return result, err //return the reference
}

// FindEnvironmentByReference returns the Environment object (interface) by reference (r)
func (c *Client) FindEnvironmentByReference(r string) (interface{}, error) {
	result, err := c.FindObject("environment", "reference", r, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Environment Reference %s not found\n", r)
	}

	return result, err //return the reference
}

// FindEnvironmentByName returns the Environment object (interface) by reference (r)
func (c *Client) FindEnvironmentByName(n string) (interface{}, error) {
	result, err := c.FindObject("environment", "name", n, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Environment Reference %s not found\n", n)
	}

	return result, err //return the reference
}

// FindEnvironmentPrimaryUser returns the reference of the environment user
func (c *Client) FindEnvironmentPrimaryUser(r string) (string, error) {
	result, err := c.FindObject("environment", "reference", r, 0)
	if err != nil {
		glog.Fatalf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		glog.V(2).Infof("Environment Reference %s not found\n", r)
	}

	userRef := result.(map[string]interface{})["primaryUser"].(string)

	return userRef, err
}
