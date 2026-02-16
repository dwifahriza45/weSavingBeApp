package database

import (
	"BE_WE_SAVING/internal/Infrastructures/config"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	DB *sqlx.DB
}

func NewStore(cfg *config.Config) *Store {
	db, err := sqlx.Connect("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("DB CONNECT ERROR:", err)
	}

	log.Println("DB CONNECTED")

	return &Store{DB: db}
}
