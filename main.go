package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mellomaths/rss-aggregator/internal/database"

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
	portString := getEnvironmentVariable("PORT")
	databaseDriver := getEnvironmentVariable("DATABASE_DRIVER")
	databaseUrl := getEnvironmentVariable("DATABASE_URL")
	conn, err := sql.Open(databaseDriver, databaseUrl)
	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	apiCfg := NewApiConfig(conn)
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
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	router.Mount("/v1", v1Router)
	server := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}
	log.Printf("Server starting on port %v", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
