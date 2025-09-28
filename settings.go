package main

import (
	"log"
	"os"
)

func getEnvironmentVariable(key string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return env
}
