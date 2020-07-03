package delphix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"
	"os"
	resty "github.com/go-resty/resty/v2"
)

// Client the structure of a client request
type Client struct {
	url string
	username string
	password string 
	restClient *resty.Client
}

// AuthSuccess holds the resty response success message
type AuthSuccess struct {
	ID, Message string
}

// RespError holds the resty response failure message
type RespError struct {
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	ErrorStruct `json:"error,omitempty"`
}

// ErrorStruct is the struct of a resty error
type ErrorStruct struct {
	Type          string `json:"type,omitempty"`
	Details       string `json:"details,omitempty"`
	ID            string `json:"id,omitempty"`
	CommandOutput string `json:"commandOutput,omitempty"`
	Diagnosis     string `json:"diagnosis,omitempty"`
}


//CreateAPIVersion returns an APISession object
func CreateAPIVersion(major int, minor int, micro int) (APIVersionStruct, error) {
	maj := new(int)
	min := new(int)
	mic := new(int)
	t := "APIVersion"
	*maj = major
	*min = minor
	*mic = micro

	apiVersion := APIVersionStruct{
		Type:  t,
		Major: maj,
		Minor: min,
		Micro: mic,
	}
	return apiVersion, nil
}

// CreateAPISession returns an APISession object
//v = APIVersion Struct
//l = Locale as an IETF BCP 47 language tag, defaults to 'en-US'.
//c = Client software identification token.
func CreateAPISession(v APIVersionStruct, l string, c string) (APISessionStruct, error) {
	if l == "" {
		l = "en-US"
	}
	if len(c) > 63 {
		err := fmt.Errorf("Client ID specified cannot be longer than 64 characters.\nYou provided %s", c)
		return APISessionStruct{}, err
	}
	apiSession := APISessionStruct{
		Type:    "APISession",
		Version: &v,
		Locale:  l,
		Client:  c,
	}
	return apiSession, nil
}

// NewClient creates a new client object
func NewClient(username, password, url string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
	}
}


// Login establishes a new client connection
func (c *Client) Login() error {

	versionStruct, err := CreateAPIVersion(1, 9, 0)
	if err != nil {
		return err
	}
	apiStruct, err := CreateAPISession(versionStruct, "", "")
	if err != nil {
		return err
	}
	// Create a Resty Client
	c.restClient = resty.New()
	c.restClient.
		SetTimeout(time.Duration(30 * time.Second)).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(apiStruct).
		Post(c.url + "/session")

	result := resp.Body()
	var resultdat map[string]interface{}
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return err
	}

	if resultdat["status"].(string) == "ERROR" {
		errorMessage := string(result)
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return err
		}
	}

	resp, err = c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetResult(AuthSuccess{}).
		SetBody(LoginRequestStruct{
			Type:     "LoginRequest",
			Username: c.username,
			Password: c.password,
		}).
		Post(c.url + "/login")
	if err != nil {
		return err
	}
	if http.StatusOK != resp.StatusCode() {
		err = fmt.Errorf("Delphix Username/Password incorrect")
		if err != nil {
			return err
		}
	}
	result = resp.Body()
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return err
	}

	if resultdat["status"].(string) == "ERROR" {
		errorMessage := string(result)
		log.Fatalf(errorMessage)
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return err
		}
	}
	
	return nil
}


func (c *Client) executeGetReturnBody(u string) (
	map[string]interface{}, error) {
	var resultMap map[string]interface{}
	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetError(&RespError{}).
		Get(c.url + u)
	if err != nil {
		return nil, err
	}

	if http.StatusOK != resp.StatusCode() { //check to make sure our query was good
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		if err != nil {
			if errorID := resp.Error().(*RespError).ErrorStruct.ID; errorID == "exception.executor.object.missing" {
				//object is missing, return that as the error
				return resultMap, fmt.Errorf(errorID)
			}
			return nil, err
		}
	}
	responseBody := resp.Body()
	// fmt.Printf("\n\nRESPONSE: %v\n\n", string(responseBody))
	if err = json.Unmarshal(responseBody, &resultMap); err != nil { //convert the json to go objects
		return nil, err
	}

	if resultMap["status"].(string) == "ERROR" {
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return nil, err
		}
	}
	return resultMap, err //return the query results
}

func (c *Client) executeListAndReturnResults(u string) (
	[]interface{}, error) {

	obj, err := c.executeGetReturnBody(u)
	if err != nil {
		return nil, err
	}
	return obj["result"].([]interface{}), err //return the query results
}


func (c *Client) executePostJobAndReturnObjectReference(u string, p interface{}) (
	interface{}, error) {

	postBody := p
	//DEBUG
	// tbEnc, err := json.Marshal(postBody)
	// fmt.Println(string(tbEnc))
	//DEBUG
	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(postBody).
		Post(c.url + u)

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

	result := resp.Body()

	var resultdat map[string]interface{}
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return nil, err
	}

	if resultdat["status"].(string) == "ERROR" {
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return nil, err
		}
	}

	reference := resultdat["result"].(string)           //grab the vdb reference
	if jobNumber, ok := resultdat["job"].(string); ok { //grab the job reference
		c.WaitforDelphixJob(jobNumber)
	}

	return reference, err
}


// WaitforDelphixJob waits for a job to complete
func (c *Client) WaitforDelphixJob(j string) error {
	var jobState string
	var err error
	for jobState != "COMPLETED" && jobState != "FAILED" && jobState != "CANCELED" {
		time.Sleep(3 * time.Second)
		resp, err := c.restClient.R().
			SetHeader("Content-Type", "application/json").
			Get(c.url + "/job/" + j)
		// explore response object
		if err != nil {
			panic(err)
		}
		s := resp.Body()

		var dat map[string]interface{}
		if err = json.Unmarshal(s, &dat); err != nil { //convert the json to go objects
			return err
		}
		results := dat["result"].(map[string]interface{}) //grab the query results
		jobState = results["jobState"].(string)
		fmt.Println(results["jobState"])
	}
	//If the job is failed or cancelled, return an error
	if jobState == "FAILED" {
		err = fmt.Errorf("Job Failed")
	} else if jobState == "CANCELED" {
		err = fmt.Errorf("Job Canceled")
	}
	return err
}

func (c *Client) executePostJobAndReturnErrOnly(u string, p interface{}) error {

	postBody := p
	//DEBUG
	// tbEnc, err := json.Marshal(postBody)
	// fmt.Println(string(tbEnc))
	//DEBUG
	resp, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(postBody).
		Post(c.url + u)

	if err != nil {
		return err
	}

	if http.StatusOK != resp.StatusCode() { //check to make sure our query was good
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return err
		}
	}

	result := resp.Body()

	var resultdat map[string]interface{}
	if err = json.Unmarshal(result, &resultdat); err != nil { //convert the json to go objects
		return err
	}

	if resultdat["status"].(string) == "ERROR" {
		errorMessage := string(resp.Body())
		err = fmt.Errorf(errorMessage)
		if err != nil {
			return err
		}
	}

	if jobNumber, ok := resultdat["job"].(string); ok { //grab the job reference
		c.WaitforDelphixJob(jobNumber)
	}

	return err
}

func (c *Client) executeReadAndReturnObject(u string) (
	interface{}, error) {
	obj, err := c.executeGetReturnBody(u)
	if err != nil {
		//return err, if anything but object does not exist
		if err.Error() != "exception.executor.object.missing" {
			return nil, err
		}
		//object is missing, return nil, not err
		return nil, nil
	}
	return obj["result"].(interface{}), err //return the query result
}

// Simple function to get Environment Variable (returns fallback string in case env variable is null)
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}