package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func New(url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return conn, nil
}
