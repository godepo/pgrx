//go:generate mockery
package example

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Repository struct {
	db DB
}

func New(db DB) *Repository {
	return &Repository{db: db}
}

const createQuery = `
INSERT INTO 
    groat_pay.payments (id, user_id, amount, kind, currency ) 
VALUES 
    ($1, $2, $3, $4, $5)`

func (repo *Repository) Create(ctx context.Context, payment Payment) error {
	_, err := repo.db.Exec(ctx,
		createQuery, payment.ID, payment.UserID, payment.Amount, payment.Kind, payment.Currency)

	return err

}

const findByIdAndUserIdQuery = `
SELECT 
    id, user_id, amount, kind, currency, state, created_at, processed_at
FROM groat_pay.payments 
WHERE user_id=$2 AND id=$1`

func (repo *Repository) FindByIDAndUserID(ctx context.Context, id uuid.UUID, userID string) (empty Payment, err error) {
	rows, err := repo.db.Query(ctx, findByIdAndUserIdQuery, id, userID)
	if err != nil {
		return empty, err
	}
	payment, err := pgx.CollectOneRow[Payment](rows, pgx.RowToStructByName[Payment])
	if err != nil {
		return empty, err
	}
	return payment, nil
}
