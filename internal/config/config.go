package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/englandrecoil/go-marketplace-service/internal/handlers"
	"github.com/joho/godotenv"
)

func Init() handlers.ApiConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".env must be created in current directory")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set ")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open connection to db: %v", err)
	}

	return handlers.ApiConfig{
		Conn: dbConn,
	}
}
