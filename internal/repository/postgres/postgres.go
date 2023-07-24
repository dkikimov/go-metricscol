package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func New(url string) (*pgx.Conn, error) {
	if len(url) == 0 {
		log.Printf("No database url provided, skipping database initialization")
		return nil, nil
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return conn, nil
}
