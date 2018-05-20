package main

import (
	"log"

	"github.com/hashicorp/packer/packer/plugin"
	"github.com/jbonachera/packer-scaleway-plugin/scaleway/volumesurrogate"
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
