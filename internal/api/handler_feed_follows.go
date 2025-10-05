package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
	"github.com/mellomaths/rss-aggregator/internal/models"
)

func (apiCfg *ApiConfig) HandleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	params := models.CreateFeedFollowParams{}
	if err := params.Decode(r); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}
	feedFollow, err := apiCfg.DATABASE.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "RECORD_CREATE_ERROR", fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, models.NewFeedFollowFromDatabase(feedFollow))
}

func (apiCfg *ApiConfig) HandleGetFeedsFollowedByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	pagination := models.PaginatedParams{}
	if err := pagination.Decode(r); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}
	if err := pagination.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	feedsFollowed, err := apiCfg.DATABASE.GetFeedsFollowedByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "RECORD_GET_ERROR", fmt.Sprintf("Error getting feeds followed by user: %v", err))
		return
	}
	data := models.NewFeedFollowsFromDatabase(feedsFollowed)
	total := len(data)
	respondWithJson(w, http.StatusOK, models.Paginated[*models.FeedFollow]{
		Data:   data,
		Total:  total,
		Offset: int(pagination.Offset),
		Limit:  int(pagination.Limit),
	})
}

func (apiCfg *ApiConfig) HandleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	params := models.DeleteFeedFollowParams{}
	if err := params.Decode(chi.URLParam(r, "feedFollowID")); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_URL_PARAMS", err.Error())
		return
	}
	err := apiCfg.DATABASE.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     params.ID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "RECORD_DELETE_ERROR", fmt.Sprintf("Error deleting feed follow: %v", err))
		return
	}
	respondWithJson(w, http.StatusNoContent, struct{}{})
}
