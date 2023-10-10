package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ResponseBuilder(w http.ResponseWriter, code int, err string, data any) {
	var success bool

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err == "" {
		success = true
	}

	json.NewEncoder(w).Encode(Response{
		Success: success,
		Error:   err,
		Data:    data,
	})
}
