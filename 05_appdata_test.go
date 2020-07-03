package delphix

import (
	"log"
	"testing"
)

func TestAppData_DeleteDatabase(t *testing.T) {
	log.Println("05 - TestAppData_DeleteDatabase")
	testDelphixAdminClient.Login()
	err := testDelphixAdminClient.DeleteAppData(appDataReference)
	if (err != nil){
		log.Fatalf("[ERROR] Delete VDB %s failed\n", appDataReference)	
		t.Errorf("[ERROR] Delete VDB %s failed\n", appDataReference)	
	}
	log.Printf("Successfully deleted VDB %s\n", appDataReference)
}