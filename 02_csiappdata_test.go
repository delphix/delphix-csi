package delphix

import (
	"log"
	"os"
	"testing"
)

var (
	appDataReference              string
	envReference                  string
	groupReference                string
	appDataSourceConfigReference  string
	appDataSourceConfigParameters interface{}
)

func TestFinders_FindGroupByname(t *testing.T) {
	log.Println("02 - TestFinders_FindGroupByname")
	groupName := os.Getenv("GROUP_NAME")
	result, err := testDelphixAdminClient.FindGroupByName(groupName)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("Group %s not found\n", groupName)
		t.Errorf("Group %s not found\n", groupName)
	}
	groupReference = result.(map[string]interface{})["reference"].(string)
	log.Printf("Group Name %s: Group Reference %s\n", groupName, groupReference)

}

func TestFinders_FindEnvironmentByName(t *testing.T) {
	log.Println("02 - TestFinders_FindEnvironmentByName")
	envName := os.Getenv("ENVIRONMENT_NAME")
	result, err := testDelphixAdminClient.FindEnvironmentByName(envName)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("Environment %s not found\n", envName)
		t.Errorf("Environment %s not found\n", envName)
	}
	envReference = result.(map[string]interface{})["reference"].(string)
	log.Printf("Environment Name %s: Environment Reference %s\n", envName, envReference)

}

func TestCSIAppData_ProvisionCloneAppData(t *testing.T) {
	fbool := new(bool)
	*fbool = false
	tbool := new(bool)
	*tbool = true
	type fields struct {
		url      string
		username string
		password string
	}
	type args struct {
		name        string
		environment string
		repo        string
		group       string
		source      string
		timestamp   string
		pvcName     string
		pvName      string
		namespace   string
		mountPath   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test1",
			want:    "A dSource reference",
			wantErr: false,
			args: args{
				name:        os.Getenv("APP_DATA_NAME"),
				repo:        os.Getenv("REPO_NAME"),
				environment: envReference,
				group:       groupReference,
				source:      os.Getenv("APP_DATA_SOURCE"),
				timestamp:   "LATEST_SNAPSHOT",
				pvcName:     os.Getenv("APP_DATA_NAME"),
				pvName:      os.Getenv("APP_DATA_NAME"),
				namespace:   os.Getenv("APP_NAMESPACE"),
				mountPath:   os.Getenv("APP_DATA_MOUNT_PATH"),
			},
		},
	}
	log.Println("02 - TestCSIAppData_ProvisionCloneAppData")
	testAppData := getEnv("DELPHIX_TEST_APPDATA", "NO")

	if testAppData == "YES" || testAppData == "yes" {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testDelphixAdminClient.Login()

				// name string, envRef string, groupRef string, sourceRef string, timestamp string, sourceConfigParameter CSISourceConfigParameterStruct
				got, err := testDelphixAdminClient.ProvisionCSIAppData(tt.args.name, tt.args.environment, tt.args.repo, tt.args.group, tt.args.source, tt.args.timestamp, tt.args.pvcName, tt.args.pvName, tt.args.namespace, tt.args.mountPath)
				if (err != nil) != tt.wantErr {
					t.Errorf("Client.ProvisionCloneAppData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got == nil {
					t.Errorf("Client.ProvisionCloneAppData() did not return a dSource reference")
					return
				}
				appDataReference = got.(string)
				log.Printf(got.(string))
			})
		}
	}
}

func TestCSIAppData_FindAppDataContainerByReference(t *testing.T) {
	log.Println("02 - TestAppData_FindAppDataContainerByReference")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindAppDataContainerByReference(appDataReference)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if (result == AppDataContainerStruct{}) {
		log.Fatalf("Database %s not found\n", appDataReference)
		t.Errorf("Database %s not found\n", appDataReference)
	}
	ref := result.Reference
	log.Printf("Database %s: Reference %s\n", appDataReference, ref)
}
