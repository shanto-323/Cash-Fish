package internal

import (
	"context"
	"database/sql"

	"settlement/pkg"

	_ "github.com/lib/pq"
)

type Repository interface {
	MakeEntry(ctx context.Context, b pkg.UserBalance) error
	GerEntry(ctx context.Context, uid string) (*pkg.UserBalance, error)
	UpdateEntry(ctx context.Context, b *pkg.UserBalance) error
}

type SettlementRepository struct {
	db *sql.DB
}

func NewRepository(dsn string) (Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &SettlementRepository{
		db: db,
	}, nil
}

func (w *SettlementRepository) MakeEntry(ctx context.Context, b pkg.UserBalance) error {
	_, err := w.db.ExecContext(
		ctx,
		`INSERT INTO balances(
			uid,
			balance,
			active
		)VALUES($1,$2,$3)`, b.UserId, b.Balance, b.Active,
	)
	return err
}

func (w *SettlementRepository) GerEntry(ctx context.Context, uid string) (*pkg.UserBalance, error) {
	balance := &pkg.UserBalance{}
	err := w.db.QueryRowContext(
		ctx,
		`SELECT 
			uid,
			balance,
			active
		FROM balances WHERE uid = $1`,
		uid,
	).Scan(
		&balance.UserId,
		&balance.Balance,
		&balance.Active,
	)
	return balance, err
}

func (w *SettlementRepository) UpdateEntry(ctx context.Context, b *pkg.UserBalance) error {
	_, err := w.db.ExecContext(
		ctx,
		`UPDATE balances SET 
			balance = $2,
			active = $3
		WHERE uid = $1`,
		b.UserId, b.Balance, b.Active,
	)
	return err
}
