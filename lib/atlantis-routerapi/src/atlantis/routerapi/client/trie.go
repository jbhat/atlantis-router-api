package client

import (
	"errors"
)

type ListTrieCommand struct {
	Info bool `short:"i" long:"info" description:"Show full info for each trie"`
}

func (c *ListTrieCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("GET", "/tries", "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)

	if c.Info && statusCode == 200 {
		err = ExpandAndPrintData("/tries/", "Tries", data)
		if err != nil {
			ErrorPrint(err)
		}
	}

	return nil
}

type GetTrieCommand struct {
	Name string `short:"n" long:"name" description:"the name of the trie"`
}

func (c *GetTrieCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}
	if c.Name == "" {
		ErrorPrint(errors.New("Please specify a trie name"))
	}

	statusCode, data, err := BuildAndSendRequest("GET", "/tries/"+c.Name, "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil

}

type AddTrieCommand struct {
	Name     string   `short:"n" long:"name" description:"the name of the trie"`
	Rules    []string `short:"r" long:"rules" description:"the list of rules for this trie"`
	Internal bool     `short:"i" long:"internal" description:"true if internal"`
}

func (c *AddTrieCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	tJsonData, err := c.GetTrieJson()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("PUT", "/tries/"+c.Name, tJsonData)
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type DeleteTrieCommand struct {
	Name string `short:"n" long:"name" description:"the name of the trie"`
}

func (c *DeleteTrieCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	if c.Name == "" {
		ErrorPrint(errors.New("Please specify a trie name"))
	}

	statusCode, data, err := BuildAndSendRequest("DELETE", "/tries/"+c.Name, "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}
