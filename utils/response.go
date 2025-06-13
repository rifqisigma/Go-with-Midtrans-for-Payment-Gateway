package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, statuscode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	json.NewEncoder(w).Encode(message)
}
func WriteError(w http.ResponseWriter, statuscode int, message string) {
	WriteJson(w, statuscode, map[string]string{
		"errpr": message,
	})
}
