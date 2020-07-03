package delphix

import (
	"log"
	"os"
	"testing"
)

var (
	databaseReference      string
	databaseGUID           string
	sourceConfigReference  string
	sourceConfigParameters interface{}
)

func TestFinders_FindDatabaseByName(t *testing.T) {
	log.Println("04 - TestFinders_FindDatabaseByName")
	databaseName := os.Getenv("APP_DATA_NAME")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindObject("database", "name", databaseName, 0)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("Database %s not found\n", databaseName)
		t.Errorf("Database %s not found\n", databaseName)
	}
	databaseReference = result.(map[string]interface{})["reference"].(string)
	log.Printf("Database %s: Referece %s\n", databaseName, databaseReference)
}

func TestFinders_FindDatabaseByReference(t *testing.T) {
	log.Println("04 - TestFinders_FindDatabaseByReference")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindObject("database", "reference", databaseReference, 0)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("Database %s not found\n", databaseReference)
		t.Errorf("Database %s not found\n", databaseReference)
	}
	databaseGUID = result.(map[string]interface{})["guid"].(string)
	log.Printf("Database %s: GUID %s\n", databaseReference, databaseGUID)
}

func TestFindSourceByContainerReference(t *testing.T) {
	log.Println("04 - TestFindSourceByContainerReference")
	dbRef := os.Getenv("TEST_DELPHIX_DB_REF")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindObject("source", "container", appDataReference, 0)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("Source %s not found\n", dbRef)
		t.Errorf("Source %s not found\n", dbRef)
	}
	sourceConfigReference = result.(map[string]interface{})["config"].(string)
	log.Printf("Database %s: SourceConfig %s\n", dbRef, sourceConfigReference)
}

func TestFindSourceConfigByReference(t *testing.T) {
	log.Println("04 - TestFindSourceConfigByReference")
	testDelphixAdminClient.Login()
	result, err := testDelphixAdminClient.FindObject("sourceconfig", "reference", sourceConfigReference, 0)
	if err != nil {
		t.Errorf("Wamp, Wamp: %s\n", err)
	}
	if result == nil {
		log.Fatalf("SourceConfig %s not found\n", sourceConfigReference)
		t.Errorf("SourceConfig %s not found\n", sourceConfigReference)
	}
	sourceConfigParameters = result.(map[string]interface{})["parameters"].(interface{})
	log.Printf("SourceConfig %s: Parameters %s\n", sourceConfigReference, sourceConfigParameters)
}
