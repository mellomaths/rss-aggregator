package api

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

type ApiConfig struct {
	DATABASE *database.Queries
}

func NewApiConfig(conn *sql.DB) *ApiConfig {
	return &ApiConfig{
		DATABASE: database.New(conn),
	}
}

func (apiCfg *ApiConfig) SetupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// v1 API routes
	v1Router := chi.NewRouter()
	// Readiness endpoint
	v1Router.Get("/healthz", handleReadiness)
	// Users endpoints
	v1Router.Post("/users", apiCfg.HandleCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandleGetUser))
	// Feeds endpoints
	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandleGetAllFeeds)
	// Feed follows endpoints
	v1Router.Post("/feeds/follows", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeedFollow))
	v1Router.Get("/feeds/follows", apiCfg.MiddlewareAuth(apiCfg.HandleGetFeedsFollowedByUser))
	v1Router.Delete("/feeds/follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandleDeleteFeedFollow))
	// Posts endpoints
	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandleGetPostsForUser))
	router.Mount("/v1", v1Router)
	return router
}
