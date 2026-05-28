package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/types"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, types.APIError{Error: message})
}
