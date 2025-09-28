package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
	"github.com/mellomaths/rss-aggregator/internal/models"
)

func (apiCfg *ApiConfig) HandleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	params := models.CreateFeedParams{}
	err := params.Decode(r)
	if err != nil {
		respondWithError(w, 400, "INVALID_REQUEST_BODY", fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}
	err = params.Validate()
	if err != nil {
		respondWithError(w, 400, "VALIDATION_ERROR", fmt.Sprintf("Error validating JSON: %v", err))
		return
	}
	feed, err := apiCfg.DATABASE.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, "RECORD_CREATE_ERROR", fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJson(w, 201, models.NewFeedFromDatabase(feed))
}
