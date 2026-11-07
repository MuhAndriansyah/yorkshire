package database

import (
	"context"

	db "github.com/MuhAndriansyah/yorkshire/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BeginTx(
	ctx context.Context,
	conn *pgxpool.Pool,
) (pgx.Tx, *db.Queries, error) {

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	qtx := db.New(conn).WithTx(tx)

	return tx, qtx, nil
}
