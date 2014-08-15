package client

import (
	"errors"
)

type ListRuleCommand struct {
	Info		bool	`short:"i" long:"info" description:"Show full info for each rule"`
}

func (c *ListRuleCommand) Execute(args []string) error {

	err := Init() 
	if err != nil {
		ErrorPrint(err)	
	}

	statusCode, data, err := BuildAndSendRequest("GET", "/rules", "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type GetRuleCommand struct {
	Name		string	`short:"n" long:"name" description:"the name of the rule"`
}

func (c *GetRuleCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}
	
	if c.Name == "" {
		ErrorPrint(errors.New("Please specify a rule name"))
	}

	statusCode, data, err := BuildAndSendRequest("GET", "rules/" + c.Name, "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type AddRuleCommand struct {
	Name            string `short:"n" long:"name" description:"the name of the rule"`
	Type     	string `short:"t" long:"type" description:"the type of the rule"`
        Value    	string `short:"v" long:"value" description:"the rule's value"`
        Next     	string `short:"x" long:"next" description:"the next ruleset"`
        Pool     	string `short:"p" long:"pool" description:"the pool to point to if this rule succeeds"`
        Internal        bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *AddRuleCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	rJson, err := c.GetRuleJson()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("PUT", "/rules/" + c.Name, rJson)
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type DeleteRuleCommand struct {
	Name            string  `short:"n" long:"name" description:"the name of the rule"`
}

func (c *DeleteRuleCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)	
	}

	if c.Name == "" {
		ErrorPrint(errors.New("Please specify a rule name"))
	}

	statusCode, data, err := BuildAndSendRequest("DELETE", "/rules/" + c.Name, "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}


