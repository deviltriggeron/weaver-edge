package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"weaver-edge/config"
)

func InitDB(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Pass, cfg.Storage.DB)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error init db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping db: %v", err)
	}

	log.Println("DB initialized successfully")
	return db, nil
}
