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


//pool with empty hosts
func testPool() cfg.Pool {
        return cfg.Pool{
                Name:     "swimming",
                Internal: true,
                Hosts: map[string]cfg.Host{},
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


func testHost() map[string]cfg.Host {
	return map[string]cfg.Host{
		"myHost": cfg.Host{"hostaddr"},
		"yourHost": cfg.Host{"uraddr"},
		"ourHost": cfg.Host{"ouraddr"},
	}
}

func testHostData() string {

	host := testHost()
	b, err := json.Marshal(host)
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


func TestGetHosts(t *testing.T){

	pool := testPool()
	
	host := testHost()
	hostData := testHostData()
	
	if err := zk.SetPool(pool);  err != nil {
		t.Fatalf("failed to put pool to test add host")
	}

	defer func() {
		if err := zk.DeletePool(pool.Name); err != nil {
			t.Fatalf("could not clean up")
		}
	}()

	if err := zk.AddHosts(pool.Name, host); err != nil {
		t.Fatalf("couldn't add hosts to attempt get")
	}

	statusCode, data, err := client.BuildAndSendRequest("GET", "/pools/" + pool.Name + "/hosts", "")
	if err != nil {
		t.Fatalf("couldn't get hosts")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect status code")
	}

	if data != hostData {
		t.Fatalf("data from get doesn't match the put data")
	}

}


func TestAddHosts(t *testing.T){

	pool := testPool()

	hostData := testHostData()

	if err := zk.SetPool(pool); err != nil {
		t.Fatalf("could not add pool to add hosts")
	}

	defer func() {
		if err := zk.DeletePool(pool.Name); err != nil {
			t.Fatalf("could not clean up")
		}
	}()

	statusCode, _, err := client.BuildAndSendRequest("PUT", "/pools/" + pool.Name + "/hosts", hostData)
	if err != nil {
		t.Fatalf("could not add hosts")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect response status from add host req")
	}


	status, data, err := client.BuildAndSendRequest("GET", "/pools/" + pool.Name + "/hosts", "")
	if err != nil {
		t.Fatalf("could not get hosts")
	}

	if status != 200 {
		t.Fatalf("incorrect status code from get hosts req")
	}

	if data != hostData {
		t.Fatalf("host data not added properly")
	}
}

func TestDeleteHosts(t *testing.T){

	pool := testPool()
	hostData := testHostData()


	
	if err := zk.SetPool(pool); err != nil {
		t.Fatalf("could not add pool to add hosts")
	}

	defer func() {
		if err := zk.DeletePool(pool.Name); err != nil {
			t.Fatalf("could not clean up")
		}
	}()

	statusCode, _, err := client.BuildAndSendRequest("PUT", "/pools/" + pool.Name + "/hosts", hostData)
	if err != nil {
		t.Fatalf("could not add hosts")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect response status from add host req")
	}
	

	hMap := make(map[string][]string, 1)
	hRay := []string{"myHost", "yourHost", "ourHost"}
	hMap["Hosts"] = hRay 

	b, err := json.Marshal(hMap)
	if err != nil {
		t.Fatalf("could not marshal hmap")
	}
	
	delHData := string(b)

	statusCode, _, err = client.BuildAndSendRequest("DELETE", "/pools/" + pool.Name + "/hosts", delHData)
	if err != nil {
		t.Fatalf("could not delete hosts")
	}	
	
	if statusCode != 200 {
		t.Fatalf("incorrect status code from delete")
	}

	hList, err := zk.GetHosts(pool.Name)

	if err != nil {
		t.Fatalf("couldnt get hosts")
	}

	if len(hList) > 0 {
		t.Fatalf("hosts not deleted properly")
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
