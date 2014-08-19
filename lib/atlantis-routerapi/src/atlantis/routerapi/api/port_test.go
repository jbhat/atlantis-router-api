package api

import (
	"fmt"
	"testing"
	"time"
	"os/exec"
	"atlantis/routerapi/client"
	"atlantis/routerapi/api"
	"atlantis/routerapi/zk"
	"encoding/json"
	cfg "atlantis/router/config"
)

const (
	DefaultAPIAddr = "8081"
	DefaultZkPort = "0.0.0.0:2181"
	DefaultUser = "kwilson"
	DefaultSecret = "pass"
)

func testPort() cfg.Port {
	return cfg.Port{
                Port:     9999,
		Trie:	  "mytrie",
                Internal: true,
	}
}

func testPortData() string {

	port := testPort()
	b, err := json.Marshal(port)
	if err != nil {
		return ""
	}

	return string(b)

}

func TestSetup(t *testing.T){

	//Start ZK server, must be using 2181 client port
	cmd := exec.Command("zkServer", "start")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not start zkServer for testing")
	}

	//configure and start the api
	err = api.Init(DefaultAPIAddr)
	if err != nil {
		t.Fatalf("failed")
	}
	zk.Init(DefaultZkPort, false)	
	//give the zk server time to start
	time.Sleep(100 * time.Millisecond)
	go api.Listen()

	client.SetDefaults("http://0.0.0.0:" + DefaultAPIAddr, DefaultUser, DefaultSecret)
}


func TestGetPort(t *testing.T){

	port := testPort()
	portData := testPortData()

	if err := zk.SetPort(port); err != nil {
		t.Fatalf("couldn't set port for get")
	}

	defer func() {
		if err := zk.DeletePort(fmt.Sprintf("%d", port.Port)); err != nil {
			t.Fatalf("couldn't clean up")
		}
	}()
	
	statusCode, data, err := client.BuildAndSendRequest("GET", "/ports/" + fmt.Sprintf("%d", port.Port), "")
	if err != nil {
		t.Fatalf("could not get port: %s", err)	
	}

	if statusCode != 200 {

		t.Fatalf("incorrect status code returned, should be 200")
	}

	if data != portData {
		t.Fatalf("Value from get not as expected \n %s \b %s", data, portData)
	}
}

func TestSetPort(t *testing.T){

	port := testPort() 
	portData := testPortData()

	statusCode, data, err := client.BuildAndSendRequest("PUT", "/ports/" + fmt.Sprintf("%d", port.Port), portData)
	if err != nil {
		t.Fatalf("Failed to send request")
	}

	if statusCode != 200 {

		t.Fatalf("Incorrect status code for response")
	}

	defer func() {
		if err := zk.DeletePort(fmt.Sprintf("%d", port.Port)); err != nil {
			t.Fatalf("Couldn't clean up port")
		}
	}()

	statusCode, data, err = client.BuildAndSendRequest("GET", "/ports/" + fmt.Sprintf("%d", port.Port), "")
	if err != nil {
		t.Fatalf("failed to send get request for set verification")
	}

	if statusCode != 200 {
		t.Fatalf("Incorrect status code for get response")
	}

	if data != portData {
		t.Fatalf("Set port failed")
	}
}

func TestDeletePort(t *testing.T){

	port := testPort()
	portData := testPortData()

	statusCode, _, err := client.BuildAndSendRequest("PUT", "/ports/" + fmt.Sprintf("%d", port.Port), portData)
	if err != nil {
		t.Fatalf("problem setting port for delete")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect set status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("DELETE", "/ports/" + fmt.Sprintf("%d", port.Port), "")
	if err != nil {
		t.Fatalf("Problem sending delete request")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect delete status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("GET", "/ports/" + fmt.Sprintf("%d", port.Port), "")
	if err != nil {
		t.Fatalf("couldn't issue get request to check if port deleted")
	}

	if statusCode != 404 {
		t.Fatalf("port not properly deleted: %d", statusCode)
	}	
}

func TestTearDown(t *testing.T){


	

	//stop zkServer
	cmd := exec.Command("zkServer", "stop")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("could not tear down zookeeper")
	}

}
