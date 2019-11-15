package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func RespondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := Message(err.Error())
	json.NewEncoder(w).Encode(response)
}

func DecodeRequest(reader io.Reader, value interface{}) error {
	err := json.NewDecoder(reader).Decode(value)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error while decoding struct %T", value))
	}

	return nil
}
