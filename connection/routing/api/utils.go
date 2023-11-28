package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func DecodeBody(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(target)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func EncodeToResponse(w http.ResponseWriter, target interface{}) {
	err := json.NewEncoder(w).Encode(target)
	if err != nil {
		log.Println("Error at encoding a response:\n" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
