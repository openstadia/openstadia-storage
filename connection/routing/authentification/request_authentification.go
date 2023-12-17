package authentification

import (
	"encoding/json"
	"fmt"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/connection/connections"
	"github.com/openstadia/openstadia-storage/connection/crud"
	"github.com/openstadia/openstadia-storage/hub_operations"
	"github.com/openstadia/openstadia-storage/models"
	"log"
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.User, *configuration.ConfigStore, *connections.ConnectionsStore)

type EnsureAuth struct {
	handler          AuthenticatedHandler
	configStore      *configuration.ConfigStore
	connectionsStore *connections.ConnectionsStore
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u, err := hub_operations.Authorize(r.Header.Get("Authorization"), ea.configStore)
	if err != nil || u == nil {
		if err != nil { // a real internal error
			log.Println(fmt.Errorf("failed to authorize a request: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
		} else { // bad credentials
			resp, _ := json.Marshal(map[string]string{"details": "Bad credentials"})
			_, err = w.Write(resp)
			if err != nil {
				log.Println(fmt.Errorf("failed to write an bad credentials authorizatioon response: %w", err))
			}
		}
		return
	}

	localUser := models.User{Id: u.Id}

	userInfo := crud.GetStorageUserInfo(&localUser)
	if !userInfo.StorageFeatureAllowed {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	localUser.UserInfo = userInfo

	ea.handler(w, r, &localUser, ea.configStore, ea.connectionsStore)
}

func NewEnsureAuth(handlerToWrap AuthenticatedHandler, store *configuration.ConfigStore, connectionsStore *connections.ConnectionsStore) *EnsureAuth {
	return &EnsureAuth{handler: handlerToWrap, configStore: store, connectionsStore: connectionsStore}
}
