package walletservice

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID             uuid.UUID `json:"id"`
	SenderId       string    `json:"sender_id"`
	ReceiverId     string    `json:"receiver_id"`
	Amount         int64     `json:"amount"`
	Note           string    `json:"note"`
	IdempotencyKey string    `json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
}

type TransactionHistoryModel struct {
	Transactions     []*TransactionModel `json:"result"`
	TotalTransection int64               `json:"total"`
	TotalPage        int64               `json:"pages"`
}

type TransactionResponseModel struct {
	PaymentId     string    `json:"payment_id"`
	PaymentStatus string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
