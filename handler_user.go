package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/auth"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

func generateRandomHex() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
}

func (apiCfg *ApiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "INVALID_REQUEST_BODY", fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}
	user, err := apiCfg.DATABASE.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		ApiKey:    generateRandomHex(),
	})
	if err != nil {
		respondWithError(w, 400, "RECORD_CREATE_ERROR", fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJson(w, 201, NewUserFromDatabase(user))
}

func (apiCfg *ApiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, 403, "FORBIDDEN", fmt.Sprintf("Authentication error: %v", err))
		return
	}
	user, err := apiCfg.DATABASE.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 401, "UNAUTHORIZED", fmt.Sprintf("Authentication error: %v", err))
		return
	}
	respondWithJson(w, 200, NewUserFromDatabase(user))
}
