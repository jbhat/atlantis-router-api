package client

import (
	"errors"
	"fmt"
)

type ListPortCommand struct {
	Info bool `short:"i" long:"info" description:"Show full info for each port"`
}

func (c *ListPortCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("GET", "/ports", "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)

	//if the info flag is set, go through and get the data for each port in
	//the list and print it's data
	if c.Info && statusCode == 200 {
		err = ExpandAndPrintData("/ports/", "Ports", data)
		if err != nil {
			ErrorPrint(err)
		}
	}

	return nil

}

type GetPortCommand struct {
	Port uint16 `short:"p" long:"port" description:"the port number"`
}

func (c *GetPortCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	if c.Port <= 0 {
		ErrorPrint(errors.New("Please specify a valid port"))
	}

	statusCode, data, err := BuildAndSendRequest("GET", fmt.Sprintf("/ports/%d", c.Port), "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type AddPortCommand struct {
	Port     uint16 `short:"p" long:"port" description:"the port number"`
	Trie     string `short:"t" long:"trie" description:"the root trie for this port"`
	Internal bool   `short:"i" long:"internal" description:"true if internal"`
}

func (c *AddPortCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	pJson, err := c.GetPortJson()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("PUT", fmt.Sprintf("/ports/%d", c.Port), pJson)
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type DeletePortCommand struct {
	Port uint16 `short:"p" long:"port" description:"the port number"`
}

func (c *DeletePortCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	if c.Port <= 0 {
		ErrorPrint(errors.New("Please specify a valid port"))
	}

	statusCode, data, err := BuildAndSendRequest("DELETE", fmt.Sprintf("/ports/%d", c.Port), "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}
