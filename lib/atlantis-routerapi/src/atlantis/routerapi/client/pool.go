package client

import (
	"errors"
)

type ListPoolCommand struct {
	Info		bool	`short:"i" long:"info" description:"Show full info for each pool"`
}

func (c *ListPoolCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("GET", "/pools", "")
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type GetPoolCommand struct {
	Name		string	`short:"n" long:"name" description:"the name of the pool"`
}

func (c *GetPoolCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}
	if c.Name == "" {
		ErrorPrint(errors.New("Please specify a name of a pool"))
	}

	statusCode, data, err := BuildAndSendRequest("GET",  "/pools/" + c.Name, "")
	if err != nil {
		ErrorPrint(err)
	} 

	Output(statusCode, data)
	return nil
}

type AddPoolCommand struct {
	Name             string   `short:"n" long:"name" description:"the name of the pool"`
        HealthCheckEvery string   `short:"e" long:"check-every" default:"5s" description:"how often to check healthz"`
        HealthzTimeout   string   `short:"z" long:"healthz-timeout" default:"5s" description:"timeout for healthz checks"`
        RequestTimeout   string   `short:"r" long:"request-timeout" default:"120s" description:"timeout for requests"`
        Status           string   `short:"s" long:"status" default:"OK" description:"the pool's status"`
        Hosts            []string `short:"H" long:"host" description:"the pool's hosts"`
        Internal         bool     `short:"i" long:"internal" description:"true if internal"`
}

func (c *AddPoolCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}		

	pJsonData, err := c.GetPoolJson()
	if err != nil {
		ErrorPrint(err)
	}

	statusCode, data, err := BuildAndSendRequest("PUT", "/pools/" + c.Name, pJsonData)
	if err != nil {
		ErrorPrint(err)
	}

	Output(statusCode, data)
	return nil
}

type DeletePoolCommand struct {
        Name            string  `short:"n" long:"name" description:"the name of the pool"`
}

func (c *DeletePoolCommand) Execute(args []string) error {

	err := Init()
	if err != nil {
		ErrorPrint(err)
	}

	if c.Name == "" {
		ErrorPrint(errors.New("Please specify a name of a pool"))
	}

	statusCode, data, err := BuildAndSendRequest("DELETE", "/pools/" + c.Name, "")
	if err != nil {
		ErrorPrint(err)
	}		
	
	Output(statusCode, data)
	return nil
}
