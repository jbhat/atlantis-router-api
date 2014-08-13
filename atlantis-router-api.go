package main


import(
	"github.com/jigish/go-flags"
	"log"
	"atlantis/routerapi/api"
	"atlantis/routerapi/zk"

)

type RouterApi struct {
	*flags.Parser
	opts 	*RouterApiOptions
	APIListenAddress	string
	ZkServerAddress		string
}

type RouterApiOptions struct {
	APIListenAddress	string `short:"A" long:"listen-address" description:"The address for the API to listen on"`
	ZkServerAddress		string `short:"Z" long:"zk-address" description:"The address of the ZK server"`
}

func NewRouterApi() *RouterApi {
	opts := &RouterApiOptions{}
	return &RouterApi{Parser: flags.NewParser(opts, flags.Default), opts: opts}
}

func (rapi *RouterApi) loadConfig(){
	rapi.Parse()

	if rapi.opts.APIListenAddress != "" {
		rapi.APIListenAddress = rapi.opts.APIListenAddress
	} else {
		rapi.APIListenAddress = "99999"
	}
	if rapi.opts.ZkServerAddress != "" {
		rapi.ZkServerAddress = rapi.opts.ZkServerAddress
	} else {
		rapi.ZkServerAdrress = "28080"
	}
}

func main() {

	rapi := NewRouterApi()
	rapi.loadConfig()
	rapi.run()
}

func (rapi *RouterApi) run() {

	api.Init(rapi.APIListenAddress)
	zk.Init(rapi.ZkServerAddress)

}
