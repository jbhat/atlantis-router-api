package client


import (
	"strings"
	"github.com/jigish/go-flags"
	"fmt"
	"os"
	"errors"
	"encoding/json"
)


var (
	APIAddress	string
	User		string
	Secret		string	
	PrettyJson 	bool
)

func SetDefaults(addr, user, secret string){

	if addr != "" {
		APIAddress = addr
	}

	if user != "" {
		User = user
	}

	if secret != "" {
		Secret = secret
	}

}


func Init() error{

	APIAddress = "0.0.0.0:99999"
	User =	"User" 
	Secret = "Pass"
	err := overlayConfig()
	if err != nil {
		return err
	}

	//if any still empty here error
	if APIAddress == "" {
		return errors.New("Please specify an APIAddress") 
	}
	if User == "" {
		return errors.New("Please specify a User name")
	}
	if Secret == "" {
		return errors.New("Please specify a secret")
	}

	return nil

}

type ClientOpts struct {
	Address string	`short:"A" long:"addr" description:"The API's address, if just a port number will be used with 0.0.0.0"`
	User	string  `short:"U" long:"user" description:"Username"`
	Secret	string	`short:"S" long:"secret" description:"Secret"`
	PrettyJson bool `short:"J" long:"pretty-json" description:"Print the returned JSON with indents"`
}

type ApiClient struct {
	*flags.Parser
}

var clientOpts = &ClientOpts{}

func New() *ApiClient {

	o := &ApiClient{flags.NewParser(clientOpts, flags.Default)}


	//Pools
	o.AddCommand("list-pools", "list the pools", "", &ListPoolCommand{})
	o.AddCommand("get-pool", "get a pool", "", &GetPoolCommand{})
	o.AddCommand("add-pool", "add a pool", "", &AddPoolCommand{})
	o.AddCommand("delete-pool", "delete a pool", "", &DeletePoolCommand{})

	//Rules
        o.AddCommand("list-rules", "list the rules", "", &ListRuleCommand{})
        o.AddCommand("get-rule", "get a rule", "", &GetRuleCommand{})
        o.AddCommand("add-rule", "add a rule", "", &AddRuleCommand{})
        o.AddCommand("delete-rule", "delete a rule", "", &DeleteRuleCommand{})

	//Tries
        o.AddCommand("list-tries", "list the tries", "", &ListTrieCommand{})
        o.AddCommand("get-trie", "get a trie", "", &GetTrieCommand{})
        o.AddCommand("add-trie", "add a trie", "", &AddTrieCommand{})
        o.AddCommand("delete-trie", "delete a trie", "", &DeleteTrieCommand{})

	//Ports
        o.AddCommand("list-ports", "list the ports", "", &ListPortCommand{})
        o.AddCommand("get-port", "get a port", "", &GetPortCommand{})
        o.AddCommand("add-port", "add a port", "", &AddPortCommand{})
        o.AddCommand("delete-port", "delete a port", "", &DeletePortCommand{})

	return o
}

func (o *ApiClient) Run(){
	o.Parse()
}

func overlayConfig() error{


	if clientOpts.Address != "" {
		if strings.Contains(clientOpts.Address, ":"){
               		 APIAddress =	clientOpts.Address 
	
       		} else {

			APIAddress = "http://0.0.0.0:" + clientOpts.Address 
		}	
	}

	if clientOpts.User != "" {
		User = clientOpts.User
	}

	if clientOpts.Secret != "" {
		Secret = clientOpts.Secret
	}

	if clientOpts.PrettyJson {
		PrettyJson = true
	}

	return nil
}

func ErrorPrint(e error){

	fmt.Printf("Failed attempting command: %s\n", e)
	os.Exit(1)

}

func Output(status int, data string){

	fmt.Printf("Status Code: %d\n\n", status)

	if PrettyJson{
		var m map[string]interface{}
		err := json.Unmarshal([]byte(data), &m)
		if err != nil {
			ErrorPrint(err)	
		}
	
		b, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			ErrorPrint(err)
		}
		
		fmt.Printf("Data: \n%s\n\n", string(b))
	
	} else {
		fmt.Printf("Data: \n%s\n\n", data)	
	}	

}
