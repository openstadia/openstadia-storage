package routing

import (
	"github.com/openstadia/openstadia-storage/connection/routing/api/storage"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func RegisterEndpoints() {
	c := cors.AllowAll()
	http.Handle("/storage/", http.StripPrefix("/storage", c.Handler(http.HandlerFunc(storage.SetUpStorageRouter().ServeHTTP))))

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
