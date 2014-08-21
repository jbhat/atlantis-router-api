package client

import (
	"testing"
	"strings"
	"atlantis/routerapi/client"
)

const (
	DefaultUser = "kwilson"
	DefaultSecret = "pass"
	DefaultAPIAddr = "http://0.0.0.0:8081"

)


func TestSetup(t *testing.T) {

	//set defaults for client package
	client.SetDefaults(DefaultAPIAddr, DefaultUser, DefaultSecret)

}

func TestBuildRequestWithData(t *testing.T){

	req, err := client.BuildRequest("GET", "/test", "non-empty")
	if err != nil {
		t.Fatalf("should not error making basic call")
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type header not set properly")
	}

	if req.Header.Get("User") != DefaultUser {
		t.Fatalf("User header not set properly")
	}

	if req.Header.Get("Secret") != DefaultSecret {
		t.Fatalf("Secret header not set properly")
	}

	if req.URL.Path !=  "/test" {
		t.Fatalf("Url/uri not built properly")
	}

	//not equality because req.URL.Host drops http/https	
	if !strings.Contains(DefaultAPIAddr, req.URL.Host) {
		t.Fatalf("Url not set properly")
	}

	if req.Method != "GET" {
		t.Fatalf("Request method not set properly")
	}	
}

func TestBuildRequestNoData(t *testing.T){

	req, err := client.BuildRequest("DELETE", "/funtime", "")
	if err != nil {
		t.Fatalf("Should not error")
	}

	//should be not set
	if req.Header.Get("Content-Type") != "" {
		t.Fatalf("Content type should not be set for dataless request")
	}


	if req.Method != "DELETE" {
		t.Fatalf("Request method not set properly")
	}		

	if req.URL.Path != "/funtime" {
		t.Fatalf("Request url/uri not set properly")
	}	

	//not equality because req.URL.Host drops the http/https	
	if !strings.Contains(DefaultAPIAddr, req.URL.Host) {
		t.Fatalf("Request url not set properly %s : %s", req.URL.Host, DefaultAPIAddr)
	}
}
