package main

import (
	"atlantis/routerapi/client"
)

func main() {
	o := client.New()
	o.Run()
}
