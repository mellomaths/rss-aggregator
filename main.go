package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mellomaths/rss-aggregator/internal/database"
	"github.com/mellomaths/rss-aggregator/internal/infra"
	"github.com/mellomaths/rss-aggregator/internal/scraper"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DATABASE *database.Queries
}

func NewApiConfig(conn *sql.DB) *ApiConfig {
	return &ApiConfig{
		DATABASE: database.New(conn),
	}
}

func main() {
	godotenv.Load(".env")
	settings := infra.NewSettings()
	conn, err := sql.Open(settings.DatabaseDriver, settings.DatabaseUrl)
	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	defer conn.Close()
	apiCfg := NewApiConfig(conn)
	rssScraper := scraper.RSSScraper{
		Database:            apiCfg.DATABASE,
		Concurrency:         10,
		TimeBetweenRequests: 10 * time.Minute,
	}
	go rssScraper.Start()
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleReadiness)
	v1Router.Post("/users", apiCfg.HandleCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandleGetUser))
	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandleGetAllFeeds)
	v1Router.Post("/feeds/follows", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeedFollow))
	v1Router.Get("/feeds/follows", apiCfg.MiddlewareAuth(apiCfg.HandleGetFeedsFollowedByUser))
	v1Router.Delete("/feeds/follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandleDeleteFeedFollow))
	router.Mount("/v1", v1Router)
	server := &http.Server{
		Addr:    ":" + settings.Port,
		Handler: router,
	}
	log.Printf("Server starting on port %v", settings.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
