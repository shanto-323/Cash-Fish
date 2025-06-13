package pkg

import (
	"time"

	"github.com/google/uuid"
)

type UserBalance struct {
	UserId  string `json:"uid"`
	Balance int64  `json:"balance"`
	Active  bool   `json:"active"`
}

type TransactionModel struct {
	ID             uuid.UUID `json:"id"`
	SenderId       string    `json:"sender_id"`
	ReceiverId     string    `json:"receiver_id"`
	Amount         int64     `json:"amount"`
	Note           string    `json:"note"`
	IdempotencyKey string    `json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
}
