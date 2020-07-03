package delphix

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	appDataGUID              string
	appDataSnapshotReference string
)

func TestAppData_FindAppDataContainerByName(t *testing.T) {
	log.Println("03 - TestAppData_FindAppDataContainerByName")
	databaseName := os.Getenv("APP_DATA_NAME")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindAppDataContainerByName(databaseName)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if (result == AppDataContainerStruct{}) {
		log.Fatalf("Database %s not found\n", databaseName)
		t.Errorf("Database %s not found\n", databaseName)
	}
	appDataReference = result.Reference
	log.Printf("Database %s: Referece %s\n", databaseName, appDataReference)
}

func TestAppData_FindAppDataContainerByReference(t *testing.T) {
	log.Println("03 - TestAppData_FindAppDataContainerByReference")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindAppDataContainerByReference(appDataReference)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if (result == AppDataContainerStruct{}) {
		log.Fatalf("Database %s not found\n", appDataReference)
		t.Errorf("Database %s not found\n", appDataReference)
	}
	appDataGUID = result.GUID
	log.Printf("Database %s: GUID %s\n", appDataReference, appDataGUID)
}

func TestAppData_FindAppDataSourceByContainerReference(t *testing.T) {
	log.Println("03 - TestAppData_FindAppDataSourceByContainerReference")
	dbRef := os.Getenv("TEST_DELPHIX_DB_REF")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindSourceByContainerReference(appDataReference)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("Database %s not found\n", dbRef)
		t.Errorf("Database %s not found\n", dbRef)
	}
	appDataSourceConfigReference = result.(map[string]interface{})["config"].(string)
	log.Printf("Database %s: SourceConfig %s\n", dbRef, appDataSourceConfigReference)
}

func TestAppData_FindCSIAppDataSourceConfigByReference(t *testing.T) {
	log.Println("03 - TestAppData_FindAppDataSourceConfigByReference")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindAppDataCSISourceConfigByReference(appDataSourceConfigReference)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if (result == CSISourceConfigStruct{}) {
		log.Fatalf("Source Config %s not found\n", appDataSourceConfigReference)
		t.Errorf("Source Config %s not found\n", appDataSourceConfigReference)
	}
	appDataSourceConfigParameters = result.Parameters
	fmt.Printf("%s\n", result.Parameters.ExportPath)
	// fmt.Printf("%+v\n", result)
	log.Printf("SourceConfig %s: Parameters %s\n", appDataSourceConfigReference, appDataSourceConfigParameters)
}

func TestAppData_SyncAppData(t *testing.T) {
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
		n string
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
			want:    "A Snapshot reference",
			wantErr: false,
			args: args{
				n: os.Getenv("APP_DATA_NAME"),
			},
		},
	}
	log.Println("03 - TestAppData_SyncAppData")
	testAppData := getEnv("DELPHIX_TEST_APPDATA", "NO")

	if testAppData == "YES" || testAppData == "yes" {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testDelphixAdminClient.Login()
				ssRef, err := testDelphixAdminClient.SyncAppData(appDataReference)
				appDataSnapshotReference = ssRef.Reference
				if (err != nil) != tt.wantErr {
					t.Errorf("Client.SyncAppData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				log.Printf("AppDataReference %s: SnapshotRepference %s\n", appDataReference, appDataSnapshotReference)
			})
		}
	}
}

func TestAppData_RollbackAppData(t *testing.T) {
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
		n string
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
			want:    "A Snapshot reference",
			wantErr: false,
			args: args{
				n: os.Getenv("APP_DATA_NAME"),
			},
		},
	}
	log.Println("03 - TestAppData_RollbackAppData")
	testAppData := getEnv("DELPHIX_TEST_APPDATA", "NO")
	if testAppData == "YES" || testAppData == "yes" {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testDelphixAdminClient.Login()
				dbObj, err := testDelphixAdminClient.FindAppDataContainerByName(tt.args.n)
				ssRef, err := testDelphixAdminClient.RollbackAppData(dbObj.Reference, appDataSnapshotReference)
				if (err != nil) != tt.wantErr {
					t.Errorf("Client.RollbackAppData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				log.Printf("Rollback to %s successfull\n", ssRef.Reference)
			})
		}
	}
}
