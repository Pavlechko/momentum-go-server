package handlers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, d any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(d)
}

func IsDecodeJSONRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return false
	}
	return true
}