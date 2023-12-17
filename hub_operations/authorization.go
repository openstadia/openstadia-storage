package hub_operations

import (
	"encoding/json"
	"fmt"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/models"
	"net/http"
)

func Authorize(auth string, store *configuration.ConfigStore) (*models.HubUser, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/me", store.GetHubSettings().HubDomain), nil)
	if err != nil {
		return nil, fmt.Errorf("failure on creating a new request: %w", err)
	}

	req.Header = http.Header{}
	req.Header.Set("Authorization", auth)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failure on creating a sending a request: %w", err)
	}

	var u models.HubUser
	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&u)
	if err != nil {
		return nil, nil
	}

	return &u, nil
}
