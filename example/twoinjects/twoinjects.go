package example

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
type Repository struct {
	firstDB  DB
	secondDB DB
}

func New(firstDB, secondDB DB) *Repository {
	return &Repository{firstDB: firstDB, secondDB: secondDB}
}
