package api

import (
	cfg "atlantis/router/config"
	"atlantis/routerapi/api"
	"atlantis/routerapi/client"
	"atlantis/routerapi/zk"
	zktest "atlantis/routerapi/zk/testutils"
	"encoding/json"
	"testing"
)

const (
	DefaultAPIAddr = "8081"
	DefaultZkPort  = 2181
	DefaultUser    = "kwilson"
	DefaultSecret  = "pass"
)

func testTrie() cfg.Trie {
	return cfg.Trie{
		Name: "breakable",
		Rules: []string{
			"myRule",
			"yourRule",
			"ourRule",
		},
		Internal: true,
	}
}

func testTrieData() string {

	trie := testTrie()
	b, err := json.Marshal(trie)
	if err != nil {
		return ""
	}

	return string(b)

}

var zkServer *zktest.ZkTestServer

func TestSetup(t *testing.T) {

	//create/start zkserver
	zkServer = zktest.NewZkTestServer(DefaultZkPort)
	if err := zkServer.Init(); err != nil {
		t.Fatalf("could not start zkServer for testing")
	}

	tmpAddr, err := zkServer.Server.Addr()
	if err != nil {
		t.Fatalf("could not get zk server addr")
	}

	//set our conn to the one created by testserver instead
	//of making a new one with zk.Init
	zk.SetZkConn(zkServer.Zk.Conn, zkServer.ZkEventChan, tmpAddr)

	//configure and start the api
	err = api.Init(DefaultAPIAddr)
	if err != nil {
		t.Fatalf("failed")
	}

	go api.Listen()

	client.SetDefaults("http://0.0.0.0:"+DefaultAPIAddr, DefaultUser, DefaultSecret)
}

func TestGetTrie(t *testing.T) {

	trie := testTrie()
	trieData := testTrieData()

	if err := zk.SetTrie(trie); err != nil {
		t.Fatalf("couldn't set trie for get")
	}

	defer func() {
		if err := zk.DeleteTrie(trie.Name); err != nil {
			t.Fatalf("couldn't clean up")
		}
	}()

	statusCode, data, err := client.BuildAndSendRequest("GET", "/tries/"+trie.Name, "")
	if err != nil {
		t.Fatalf("could not get trie: %s", err)
	}

	if statusCode != 200 {

		t.Fatalf("incorrect status code returned, should be 200")
	}

	if data != trieData {
		t.Fatalf("Value from get not as expected \n %s \b %s", data, trieData)
	}
}

func TestSetTrie(t *testing.T) {

	trie := testTrie()
	trieData := testTrieData()

	statusCode, data, err := client.BuildAndSendRequest("PUT", "/tries/"+trie.Name, trieData)
	if err != nil {
		t.Fatalf("Failed to send request")
	}

	if statusCode != 200 {

		t.Fatalf("Incorrect status code for response")
	}

	defer func() {
		if err := zk.DeleteTrie(trie.Name); err != nil {
			t.Fatalf("Couldn't clean up trie")
		}
	}()

	statusCode, data, err = client.BuildAndSendRequest("GET", "/tries/"+trie.Name, "")
	if err != nil {
		t.Fatalf("failed to send get request for set verification")
	}

	if statusCode != 200 {
		t.Fatalf("Incorrect status code for get response")
	}

	if data != trieData {
		t.Fatalf("Set trie failed")
	}
}

func TestDeleteTrie(t *testing.T) {

	trie := testTrie()
	trieData := testTrieData()

	statusCode, _, err := client.BuildAndSendRequest("PUT", "/tries/"+trie.Name, trieData)
	if err != nil {
		t.Fatalf("problem setting trie for delete")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect set status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("DELETE", "/tries/"+trie.Name, "")
	if err != nil {
		t.Fatalf("Problem sending delete request")
	}

	if statusCode != 200 {
		t.Fatalf("incorrect delete status code")
	}

	statusCode, _, err = client.BuildAndSendRequest("GET", "/tries/"+trie.Name, "")
	if err != nil {
		t.Fatalf("couldn't issue get request to check if trie deleted")
	}

	if statusCode != 404 {
		t.Fatalf("trie not properly deleted: %d", statusCode)
	}
}

func TestTearDown(t *testing.T) {

	zk.KillConnection()
	if err := zkServer.Destroy(); err != nil {
		t.Fatalf("could not destory zk server")
	}
}
