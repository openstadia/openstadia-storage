package authentification

import (
	"encoding/json"
	"github.com/openstadia/openstadia-storage/connection/crud"
	"github.com/openstadia/openstadia-storage/hub_operations"
	"github.com/openstadia/openstadia-storage/models"
	"log"
	"net/http"
)

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.HubUser)

type EnsureAuth struct {
	handler AuthenticatedHandler
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := hub_operations.Authorize(r.Header.Get("Authorization"))
	if err != nil {
		log.Println("An error on authorizing:\n" + err.Error())
	}
	userInfo := crud.GetStorageUserInfo(u)
	if !userInfo.StorageFeatureAllowed {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err != nil || u == nil {
		var details string
		if err == nil {
			details = "Bad credentials"
		} else {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, _ := json.Marshal(map[string]string{"details": details})
		_, err := w.Write(resp)
		if err != nil {
			log.Fatalln(err.Error())
		}
		return
	}

	ea.handler(w, r, u)
}

func NewEnsureAuth(handlerToWrap AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}
