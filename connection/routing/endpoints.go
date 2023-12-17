package routing

import (
	"fmt"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/routing/api/storage"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func RegisterEndpoints(configStore *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) {
	c := cors.AllowAll()
	http.Handle("/storage/", http.StripPrefix("/storage", c.Handler(http.HandlerFunc(storage.SetUpStorageRouter(configStore, connectionsStore).ServeHTTP))))

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to register endpoints: %w", err))
	}
}
