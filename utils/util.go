package utils

import (
	"encoding/json"
	"net/http"
)

func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func RespondSuccess(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func RespondError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := Message(err.Error())
	json.NewEncoder(w).Encode(response)
}
