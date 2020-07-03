package delphix

import (
	"fmt"
	"os"
	"testing"
	"log"
)

var (
	testDelphixAdminClient Client
)

func init() {	
	testDelphixAdminClient.url = os.Getenv("TEST_DELPHIX_URL") + "/resources/json/delphix"
	testDelphixAdminClient.username = os.Getenv("TEST_DELPHIX_USERNAME")
	testDelphixAdminClient.password = os.Getenv("TEST_DELPHIX_PASSWORD")
}

func TestClient_Login(t *testing.T) {
	log.Println("01 - TestClient_Login")
	fmt.Println("Client: ", testDelphixAdminClient.url, testDelphixAdminClient.username, testDelphixAdminClient.password)
	err := testDelphixAdminClient.Login()
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}
}
