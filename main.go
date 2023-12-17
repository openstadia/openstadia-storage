package main

import (
	"fmt"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/routing"
	"log"
)

func main() {
	configStore, err := configuration.Load()
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to load the config: %w", err))
	}

	connectionsStore, err := connections.InitAllConnections(configStore)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to init connections: %w", err))
	}

	routing.RegisterEndpoints(configStore, connectionsStore) // blocking call
}
