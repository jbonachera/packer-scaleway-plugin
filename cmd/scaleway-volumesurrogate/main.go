package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/jbonachera/packer-scaleway-plugin/scaleway/volumesurrogate"
	"log"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		log.Fatal(err)
	}
	b := &volumesurrogate.Builder{}
	server.RegisterBuilder(b)
	server.Serve()
}
