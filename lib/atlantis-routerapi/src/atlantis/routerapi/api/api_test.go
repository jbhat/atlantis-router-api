package api

import (
	"testing"
	"io/ioutil"
	"net/http"
	"atlantis/routerapi/client"
	"atlantis/routerapi/api"
	"atlantis/routerapi/zk"
	zktest "atlantis/routerapi/zk/testutils"
)

const (
        DefaultAPIAddr = "8081"
        DefaultZkPort = 2181
        DefaultUser = "kwilson"
        DefaultSecret = "pass"
)

var zkServer *zktest.ZkTestServer

func TestSetup(t *testing.T){

	//create/start the zkserver and set the conn in the zk package
        zkServer = zktest.NewZkTestServer(DefaultZkPort)
        if err := zkServer.Init(); err != nil {
                t.Fatalf("could not start zkServer for testing")
        }

        tmpAddr, err := zkServer.Server.Addr()
        if err != nil {
                t.Fatalf("could not get zk server addr")
        }

        //set our connection in the zk package to use the one made by zkserver
        //instead of creating an entire new one with zk.Init
        zk.SetZkConn(zkServer.Zk.Conn, zkServer.ZkEventChan, tmpAddr)

	
        //configure and start the api
        err = api.Init(DefaultAPIAddr)
        if err != nil {
                t.Fatalf("failed")
        }

        go api.Listen()

        client.SetDefaults("http://0.0.0.0:" + DefaultAPIAddr, DefaultUser, DefaultSecret)
}

func TestHealthz(t *testing.T){

	req, err := client.BuildRequest("GET", "/healthz", "")
	if err != nil {
		t.Fatalf("Error making healthz request")
	}

	c := &http.Client{}
	
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Error sending healthz request %s\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Incorrect status code from healthz")
	}

	if resp.Header.Get("Server-Status") != "OK" {
		t.Fatalf("Healthz not returning ok server status: %s", resp.Header.Get("Server-Status"))
	}

}

func TestNotFound(t *testing.T){

	req, err := client.BuildRequest("GET", "/nonexistentpage", "")
	if err != nil {
		t.Fatalf("Error building not found request")
	}

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Error sending not found request")
	}

	if resp.StatusCode != 404 {
		t.Fatalf("Incorrect status code from not found request")
	}	

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read not found response body")
	}

	if string(body) != "404 Not Found" {
		t.Fatalf("Not found body not set correctly")
	}

}


func TestTearDown(t *testing.T){

	zk.KillConnection()
        if err := zkServer.Destroy(); err != nil {
                t.Fatalf("error destroying zookeeper")
        }
}
