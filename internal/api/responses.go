package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON payload %v: %v", payload, err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, status int, errorCode string, errorMessage string) {
	log.Printf("Response with %v error: %v: %v", status, errorCode, errorMessage)
	type ErrorResponse struct {
		ErrorCode    string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}
	respondWithJson(w, status, ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	})
}
