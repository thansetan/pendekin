package helper

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func ResponseBuilder[T any](w http.ResponseWriter, code int, err string, data T) {
	var success bool

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err == "" {
		success = true
	}

	json.NewEncoder(w).Encode(Response[T]{
		Success: success,
		Error:   err,
		Data:    data,
	})
}
