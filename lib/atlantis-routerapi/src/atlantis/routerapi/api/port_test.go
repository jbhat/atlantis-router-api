package api

import (
	"fmt"
	"testing"
	"atlantis/routerapi/client"
	"atlantis/routerapi/api"
	"atlantis/routerapi/zk"
	zktest "atlantis/routerapi/zk/testutils"
	"encoding/json"
	cfg "atlantis/router/config"
)

const (
	DefaultAPIAddr = "8081"
	DefaultZkPort = 2181
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

var zkServer *zktest.ZkTestServer

func TestSetup(t *testing.T){

	//create/start zk server
	zkServer = zktest.NewZkTestServer(DefaultZkPort)
	if err := zkServer.Init(); err != nil {
		t.Fatalf("Could not start zk server")
	}

	tmpAddr, err := zkServer.Server.Addr()
	if err != nil {
		t.Fatalf("Could not get zk server addr")
	} 
	//set connection to use one created by server instead of making new one	
	zk.SetZkConn(zkServer.Zk.Conn, zkServer.ZkEventChan, tmpAddr)

	//configure and start the api
	err = api.Init(DefaultAPIAddr)
	if err != nil {
		t.Fatalf("failed")
	}
	
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

	zk.KillConnection()
	if err := zkServer.Destroy(); err != nil {
		t.Fatalf("error destroying zookeeper")
	}
}
