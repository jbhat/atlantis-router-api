package api

import (
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

func testRule() cfg.Rule {
	return cfg.Rule{
                Name:     "breakable",
		Type:	  "unenforceable",
		Value:	  "worthless",
		Next:	  "nextRule",
		Pool:	  "pool",
                Internal: true,
	}
}

func testRuleData() string {

	rule := testRule()
	b, err := json.Marshal(rule)
	if err != nil {
		return ""
	}

	return string(b)

}

var zkServer *zktest.ZkTestServer

func TestSetup(t *testing.T){

	//create/start the zkserver 
	zkServer = zktest.NewZkTestServer(DefaultZkPort)
	if err := zkServer.Init(); err != nil {
		t.Fatalf("could not start zk server")
	}

	tmpAddr, err := zkServer.Server.Addr()
	if err != nil {
		t.Fatalf("could not get server addr")
	}

	//set conn to use one created by server instead of making a new one 
	//with zk.Init	
	zk.SetZkConn(zkServer.Zk.Conn, zkServer.ZkEventChan, tmpAddr)


			
	//configure and start the api
	err = api.Init(DefaultAPIAddr)
	if err != nil {
		t.Fatalf("failed")
	}

	go api.Listen()

	client.SetDefaults("http://0.0.0.0:" + DefaultAPIAddr, DefaultUser, DefaultSecret)
}


func TestGetRule(t *testing.T){

	rule := testRule()
	ruleData := testRuleData()

	if err := zk.SetRule(rule); err != nil {
		t.Fatalf("couldn't set rule for get")
	}

	defer func() {
		if err := zk.DeleteRule(rule.Name); err != nil {
			t.Fatalf("couldn't clean up")
		}
	}()
	
	statusCode, data, err := client.BuildAndSendRequest("GET", "/rules/" + rule.Name, "")
	if err != nil {
		t.Fatalf("could not get rule: %s", err)	
	}

	if statusCode != 200 {

		t.Fatalf("incorrect status code returned, should be 200")
	}

	if data != ruleData {
		t.Fatalf("Value from get not as expected \n %s \b %s", data, ruleData)
	}
}

func TestSetRule(t *testing.T){

	rule := testRule() 
	ruleData := testRuleData()

	statusCode, data, err := client.BuildAndSendRequest("PUT", "/rules/" + rule.Name, ruleData)
	if err != nil {
		t.Fatalf("Failed to send request")
	}

	if statusCode != 200 {

		t.Fatalf("Incorrect status code for response")
	}

	defer func() {
		if err := zk.DeleteRule(rule.Name); err != nil {
			t.Fatalf("Couldn't clean up rule")
		}
	}()

	statusCode, data, err = client.BuildAndSendRequest("GET", "/rules/" + rule.Name, "")
	if err != nil {
		t.Fatalf("failed to send get request for set verification")
	}

	if statusCode != 200 {
		t.Fatalf("Incorrect status code for get response")
	}

	if data != ruleData {
		t.Fatalf("Set rule failed")
	}
}

func TestDeleteRule(t *testing.T){

	rule := testRule()
	ruleData := testRuleData()

	statusCode, _, err := client.BuildAndSendRequest("PUT", "/rules/" + rule.Name, ruleData)
	if err != nil {
		t.Fatalf("problem setting rule for delete")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect set status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("DELETE", "/rules/" + rule.Name, "")
	if err != nil {
		t.Fatalf("Problem sending delete request")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect delete status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("GET", "/rules/" + rule.Name, "")
	if err != nil {
		t.Fatalf("couldn't issue get request to check if rule deleted")
	}

	if statusCode != 404 {
		t.Fatalf("rule not properly deleted: %d", statusCode)
	}	
}

func TestTearDown(t *testing.T){

	zk.KillConnection()
	if err := zkServer.Destroy(); err != nil {
		t.Fatalf("could not destroy zookeeper")
	}	
}
