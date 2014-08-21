package main

import (
	"atlantis/routerapi/api"
	"atlantis/routerapi/zk"
	"github.com/jigish/go-flags"
)

type RouterApi struct {
	*flags.Parser
	opts             *RouterApiOptions
	APIListenAddress string
	ZkServerAddress  string
	ZkNeedConfig     bool
}

type RouterApiOptions struct {
	APIListenAddress string `short:"A" long:"listen-address" description:"The address for the API to listen on"`
	ZkServerAddress  string `short:"Z" long:"zk-address" description:"The address of the ZK server"`
	InitZK           string `short:"I" long:"init-zk" description:"Initialize zk?"`
}

func NewRouterApi() *RouterApi {
	opts := &RouterApiOptions{}
	return &RouterApi{Parser: flags.NewParser(opts, flags.Default), opts: opts}
}

func (rapi *RouterApi) setConfig(options RouterApiOptions) {

	rapi.APIListenAddress = options.APIListenAddress
	rapi.ZkServerAddress = options.ZkServerAddress
	rapi.ZkNeedConfig = (options.InitZK != "")

}
func (rapi *RouterApi) loadConfig() {
	rapi.Parse()

	if rapi.opts.APIListenAddress != "" {
		rapi.APIListenAddress = rapi.opts.APIListenAddress
	} else {
		rapi.APIListenAddress = "99999"
	}
	if rapi.opts.ZkServerAddress != "" {
		rapi.ZkServerAddress = rapi.opts.ZkServerAddress
	} else {
		rapi.ZkServerAddress = "28080"
	}
	if rapi.opts.InitZK == "y" {
		rapi.ZkNeedConfig = true
	}
}

func main() {

	rapi := NewRouterApi()
	rapi.loadConfig()
	rapi.run()
}

func (rapi *RouterApi) run() {

	api.Init(rapi.APIListenAddress)
	zk.Init(rapi.ZkServerAddress, rapi.ZkNeedConfig)

	api.Listen()

}
