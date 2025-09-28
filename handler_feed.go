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
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}
	err = params.Validate()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
		respondWithError(w, http.StatusBadRequest, "RECORD_CREATE_ERROR", fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, models.NewFeedFromDatabase(feed))
}

func (apiCfg *ApiConfig) HandleGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	pagination := models.PaginatedParams{}
	if err := pagination.Decode(r); err != nil {
		respondWithError(w, http.StatusBadRequest, "PAGINATION_ERROR", fmt.Sprintf("Error getting feeds: %v", err))
		return
	}
	feeds, err := apiCfg.DATABASE.GetAllFeeds(r.Context(), database.GetAllFeedsParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "RECORD_GET_ERROR", fmt.Sprintf("Error getting feeds: %v", err))
		return
	}
	data := models.NewFeedsFromDatabase(feeds)
	respondWithJson(w, http.StatusOK, models.Paginated[*models.Feed]{
		Data:   data,
		Total:  len(data),
		Offset: int(pagination.Offset),
		Limit:  int(pagination.Limit),
	})
}
