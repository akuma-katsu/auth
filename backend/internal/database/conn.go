package database

import (
	"auth/backend/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg *config.AppConfig) *Storage {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DataBaseUsername, cfg.DataBasePassword, cfg.DataBaseAddress, cfg.DataBasePort, cfg.DatabaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}

	return &Storage{DB: db}
}
