package api

import (
	"fmt"
	"net/http"

	"github.com/mellomaths/rss-aggregator/internal/database"
	"github.com/mellomaths/rss-aggregator/internal/models"
)

func (apiCfg *ApiConfig) HandleGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	paginated := models.PaginatedParams{}
	if err := paginated.Decode(r); err != nil {
		respondWithError(w, http.StatusBadRequest, "PAGINATION_ERROR", fmt.Sprintf("Error getting posts: %v", err))
		return
	}
	posts, err := apiCfg.DATABASE.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  paginated.Limit,
		Offset: paginated.Offset,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "RECORD_GET_ERROR", fmt.Sprintf("Error getting posts: %v", err))
		return
	}
	data := models.NewPostsFromDatabase(posts)
	respondWithJson(w, http.StatusOK, models.Paginated[*models.Post]{
		Data:   data,
		Total:  len(data),
		Offset: int(paginated.Offset),
		Limit:  int(paginated.Limit),
	})
}
