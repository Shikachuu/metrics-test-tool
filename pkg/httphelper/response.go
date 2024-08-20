package httphelper

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func WriteJSONResponse[T any](w http.ResponseWriter, status int, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
