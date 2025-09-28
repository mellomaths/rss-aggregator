package infra

import (
	"log"
	"os"
)

type Settings struct {
	Port           string
	DatabaseDriver string
	DatabaseUrl    string
}

func NewSettings() *Settings {
	port := getEnvironmentVariable("PORT")
	databaseDriver := getEnvironmentVariable("DATABASE_DRIVER")
	databaseUrl := getEnvironmentVariable("DATABASE_URL")
	return &Settings{
		Port:           port,
		DatabaseDriver: databaseDriver,
		DatabaseUrl:    databaseUrl,
	}
}

func getEnvironmentVariable(key string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return env
}
