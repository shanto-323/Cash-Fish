package walletservice

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	MakeTransaction(ctx context.Context, wallet TransactionModel) error
	TransactionsStatus(ctx context.Context, payment_id string) (*TransactionModel, error)
	TransactionsHistory(ctx context.Context, id string, limit int64, offset int64) ([]*TransactionModel, error)
	TotalTransaction(ctx context.Context, id string) (*int64, error)
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(dsn string) (Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &TransactionRepository{
		db: db,
	}, nil
}

func (w *TransactionRepository) MakeTransaction(ctx context.Context, wallet TransactionModel) error {
	_, err := w.db.ExecContext(
		ctx,
		`INSERT INTO payments(
			id,
			sender_id,
			receiver_id,
			amount,
			note,
			imenpotency_key,
			created_at,
		)VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		wallet.ID, wallet.SenderId, wallet.ReceiverId, wallet.Amount, wallet.Note, wallet.IdempotencyKey, wallet.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (w *TransactionRepository) TransactionsStatus(ctx context.Context, payment_id string) (*TransactionModel, error) {
	transaction := &TransactionModel{}
	err := w.db.QueryRowContext(
		ctx,
		`SELECT 
			id,
			sender_id,
			receiver_id,
			amount,
			note,
			imenpotency_key,
			created_at
		FROM payments WHERE id = $1`,
		payment_id,
	).Scan(
		transaction.ID,
		transaction.SenderId,
		transaction.ReceiverId,
		transaction.Amount,
		transaction.Note,
		transaction.IdempotencyKey,
		transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (w *TransactionRepository) TransactionsHistory(ctx context.Context, id string, limit int64, offset int64) ([]*TransactionModel, error) {
	rows, err := w.db.QueryContext(
		ctx,
		`SELECT 
			id,
			sender_id,
			receiver_id,
			amount,
			note,
			imenpotency_key,
			created_at
		FROM payments 
		WHERE sender_id = $1 OR receiver_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		id, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	transactions := []*TransactionModel{}
	for rows.Next() {
		transaction := &TransactionModel{}
		rows.Scan(
			transaction.ID,
			transaction.SenderId,
			transaction.ReceiverId,
			transaction.Amount,
			transaction.Note,
			transaction.IdempotencyKey,
			transaction.CreatedAt,
		)
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (w *TransactionRepository) TotalTransaction(ctx context.Context, id string) (*int64, error) {
	var count int64
	err := w.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
		FROM payments 
		WHERE sender_id = $1 OR receiver_id = $1`,
	).Scan(
		&count,
	)
	if err != nil {
		return nil, err
	}
	return &count, nil
}
