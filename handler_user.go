package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

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
	})
	if err != nil {
		respondWithError(w, 400, "RECORD_CREATE_ERROR", fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJson(w, 201, NewUserFromDatabase(user))
}
