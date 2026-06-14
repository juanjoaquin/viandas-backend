package database

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/juanjoaquin/viandas-backend/settings"
)

func New(ctx context.Context, s *settings.Settings) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.DB.Host, s.DB.Port, s.DB.User, s.DB.Password, s.DB.Name,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	return db, nil
}
