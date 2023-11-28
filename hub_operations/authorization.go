package hub_operations

import (
	"encoding/json"
	"errors"
	"github.com/openstadia/openstadia-storage/configuration"
	"github.com/openstadia/openstadia-storage/models"
	"net/http"
)

func Authorize(auth string) (*models.HubUser, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", configuration.Settings.HubDomain+"/users/me", nil)
	if err != nil {
		return nil, errors.New("Error on http.NewRequest():\n" + err.Error())
	}

	req.Header = http.Header{}
	req.Header.Set("Authorization", auth)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error on client.Do(request):\n" + err.Error())
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
