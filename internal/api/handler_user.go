package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
	"github.com/mellomaths/rss-aggregator/internal/models"
)

func generateRandomHex() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	hash := sha256.Sum256(randomBytes)
	return hex.EncodeToString(hash[:])
}

func (apiCfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := models.CreateUserParams{}
	err := params.Decode(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}
	err = params.Validate()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
		respondWithError(w, http.StatusBadRequest, "RECORD_CREATE_ERROR", fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, models.NewUserFromDatabase(user))
}

func (apiCfg *ApiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, models.NewUserFromDatabase(user))
}
