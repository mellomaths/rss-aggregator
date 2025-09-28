package main

import (
	"fmt"
	"net/http"

	"github.com/mellomaths/rss-aggregator/internal/auth"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

type AuthenticatedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(handler AuthenticatedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, "AUTHENTICATION_ERROR", fmt.Sprintf("Authentication error: %v", err))
			return
		}
		user, err := apiCfg.DATABASE.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 401, "AUTHENTICATION_ERROR", fmt.Sprintf("Authentication error: %v", err))
			return
		}
		handler(w, r, user)
	}
}
