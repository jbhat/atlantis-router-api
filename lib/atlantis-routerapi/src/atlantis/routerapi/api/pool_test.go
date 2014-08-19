package api

import (
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

func testPool() cfg.Pool {
	return cfg.Pool{
                Name:     "swimming",
                Internal: true,
                Hosts: map[string]cfg.Host{
                        "test0": cfg.Host{
                                Address: "localhost:8080",
                        },
                        "test1": cfg.Host{
                                Address: "localhost:8081",
                        },
                },
                Config: cfg.PoolConfig{
                        HealthzEvery:   "0s",
                        HealthzTimeout: "0s",
                        RequestTimeout: "0s",
                        Status:         "ITSCOMPLICATED",
                },
        }

}

func testPoolData() string {

	pool := testPool()
	b, err := json.Marshal(pool)
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


func TestGetPool(t *testing.T){

	pool := testPool()
	poolData := testPoolData()

	if err := zk.SetPool(pool); err != nil {
		t.Fatalf("couldn't set pool for get")
	}
	if err := zk.AddHosts(pool.Name, pool.Hosts); err != nil {
		t.Fatalf("couldn't set pool for get")
	}

	defer func() {
		if err := zk.DeletePool(pool.Name); err != nil {
			t.Fatalf("couldn't clean up")
		}
	}()
	
	statusCode, data, err := client.BuildAndSendRequest("GET", "/pools/" + pool.Name, "")
	if err != nil {
		t.Fatalf("could not get pool: %s", err)	
	}

	if statusCode != 200 {

		t.Fatalf("incorrect status code returned, should be 200")
	}

	if data != poolData {
		t.Fatalf("Value from get not as expected \n %s \b %s", data, poolData)
	}
}

func TestSetPool(t *testing.T){

	pool := testPool() 
	poolData := testPoolData()

	statusCode, data, err := client.BuildAndSendRequest("PUT", "/pools/" + pool.Name, poolData)
	if err != nil {
		t.Fatalf("Failed to send request")
	}

	if statusCode != 200 {

		t.Fatalf("Incorrect status code for response")
	}

	defer func() {
		if err := zk.DeletePool(pool.Name); err != nil {
			t.Fatalf("Couldn't clean up pool")
		}
	}()

	statusCode, data, err = client.BuildAndSendRequest("GET", "/pools/" + pool.Name, "")
	if err != nil {
		t.Fatalf("failed to send get request for set verification")
	}

	if statusCode != 200 {
		t.Fatalf("Incorrect status code for get response")
	}

	if data != poolData {
		t.Fatalf("Set pool failed")
	}
}

func TestDeletePool(t *testing.T){

	pool := testPool()
	poolData := testPoolData()

	statusCode, _, err := client.BuildAndSendRequest("PUT", "/pools/" + pool.Name, poolData)
	if err != nil {
		t.Fatalf("problem setting pool for delete")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect set status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("DELETE", "/pools/" + pool.Name, "")
	if err != nil {
		t.Fatalf("Problem sending delete request")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect delete status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("GET", "/pools/" + pool.Name, "")
	if err != nil {
		t.Fatalf("couldn't issue get request to check if pool deleted")
	}

	if statusCode != 404 {
		t.Fatalf("pool not properly deleted: %d", statusCode)
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
