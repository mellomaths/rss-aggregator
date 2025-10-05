package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/mellomaths/rss-aggregator/internal/api"
	"github.com/mellomaths/rss-aggregator/internal/infra"
	"github.com/mellomaths/rss-aggregator/internal/scraper"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	settings := infra.NewSettings()
	conn, err := sql.Open(settings.DatabaseDriver, settings.DatabaseUrl)
	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	defer conn.Close()
	apiCfg := api.NewApiConfig(conn)
	rssScraper := scraper.RSSScraper{
		Database:            apiCfg.DATABASE,
		Concurrency:         10,
		TimeBetweenRequests: 10 * time.Minute,
	}
	go rssScraper.Start()
	router := apiCfg.SetupRouter()
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
