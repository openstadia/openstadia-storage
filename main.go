package main

import (
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/routing"
	"log"
)

func main() {
	if !configuration.Init() {
		// some important settings are not set, no reason to continue
		return
	}
	err := connections.InitMinioClients()
	if err != nil {
		log.Fatalln(err)
	}
	//err = connections.StartBoltDb()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	routing.RegisterEndpoints() // blocking call
}
