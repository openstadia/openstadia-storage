package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func DecodeBody(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(target)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return false
	}
	return true
}

func EncodeToResponse(w http.ResponseWriter, target interface{}) {
	encoded, err := json.Marshal(target)
	if err != nil {
		log.Println(fmt.Errorf("failed to encode a response: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(encoded)
	if err != nil {
		log.Println(fmt.Errorf("failed to write bytes to a response: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
