/*
Package main - entrypoint for Delphix CSI Driver
Mainteiner: Daniel Stolf
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/DSTOLF/delphix-go-sdk"
	delphixdriver "github.com/delphix/delphix-csi/pkg/delphixplugin"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	endpoint        = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName      = flag.String("drivername", "defs.csi.delphix.com", "name of the driver")
	nodeID          = flag.String("nodeid", "", "node id")
	envName         = flag.String("envName", "", "Environment Name in Delphix")
	repositoryName  = flag.String("repositoryName", "", "Repository Name in Delphix used to provision the PVs")
	sourceRef       = flag.String("sourceRef", "", "Soruce Reference for Empty Volume in Delphix in Delphix, that will be cloned to provision PVs")
	groupName       = flag.String("groupName", "", "Group Name in Delphix where PVs will be provisioned")
	mountPath       = flag.String("mountPath", "", "Mount Path  where PVs will be initially mounted in Bastion Environment")
	delphixEndpoint = flag.String("url", "", "Delphix Engine URL")
	username        = flag.String("username", "", "Delphix Username")
	password        = flag.String("password", "", "Delphix Password")
	showVersion     = flag.Bool("version", false, "Show version.")
	// Set by the build process
	version = ""
)

func main() {
	flag.Parse()

	if *showVersion {
		baseName := path.Base(os.Args[0])
		fmt.Println(baseName, version)
		return
	}
	handle()
	os.Exit(0)
}

func handle() {
	var delphixURL = *delphixEndpoint + "/resources/json/delphix"

	var delphixClient = delphix.NewClient(*username, *password, delphixURL)
	err := delphixClient.Login()
	if err != nil {
		fmt.Println("Error connecting to Delphix Engine: ", err.Error())
		os.Exit(1)
	}

	// func NewDriver                     (driverName, nodeID,   endpoint, delphixClient, environmentName, groupName, repName, sourceRef, mPath, version string) (*driver, error) {
	driver, err := delphixdriver.NewDriver(*driverName, *nodeID, *endpoint, *delphixClient, *envName, *groupName, *repositoryName, *sourceRef, *mountPath, version)
	if err != nil {
		fmt.Println("Failed to initialize driver: ", err.Error())
		os.Exit(1)
	}
	driver.Run()
}
