package walletservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo        Repository
	redisClient *redis.Client
	publisher   *amqp.Channel
}

func NewService(repo Repository, redisClient *redis.Client, publisher *amqp.Channel) Service {
	return Service{
		repo:        repo,
		redisClient: redisClient,
		publisher:   publisher,
	}
}

func (s *Service) CreateNewTransection(ctx context.Context, senderId, receiverId, note, idempotencyKey string, amount float64) (*TransactionModel, error) {
	catchKey := fmt.Sprintf("IdempotencyKey:%s", idempotencyKey)
	exits, err := s.redisClient.Get(ctx, catchKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if exits != "" {
		return nil, fmt.Errorf("one request on process")
	}

	transactionModel := TransactionModel{
		ID:             uuid.New(),
		SenderId:       senderId,
		ReceiverId:     receiverId,
		Note:           note,
		IdempotencyKey: catchKey,
		Amount:         int64(amount * 100),
		CreatedAt:      time.Now(),
	}

	err = s.repo.MakeTransaction(ctx, transactionModel)
	if err != nil {
		return nil, err
	}

	if err := s.redisClient.Set(ctx, catchKey, "true", 30*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("failed to set-up idempotency key: %s", err)
	}

	if err := s.publishPaymentEvent(ctx, transactionModel); err != nil {
		return nil, fmt.Errorf("failed to publish payment event: %s", err)
	}

	return &transactionModel, nil
}

func (s *Service) GetTransection(ctx context.Context, payment_id string) (*TransactionModel, error) {
	return s.repo.TransactionsStatus(ctx, payment_id)
}

func (s *Service) GetTransectionHistory(ctx context.Context, id string, limit, offset int64) (*TransactionHistoryModel, error) {
	if limit < 10 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	type result struct {
		transactions []*TransactionModel
		total        int64
		err          error
	}

	ch := make(chan result, 2)
	go func() {
		t, err := s.repo.TransactionsHistory(ctx, id, limit, offset)
		select {
		case ch <- result{transactions: t, err: err}:
		case <-ctx.Done():
			return
		}
	}()

	go func() {
		t, err := s.repo.TotalTransaction(ctx, id)
		select {
		case ch <- result{total: *t, err: err}:
		case <-ctx.Done():
			return
		}
	}()

	var transections []*TransactionModel
	var total int64
	var err error

	for i := 0; i < 2; i++ {
		res := <-ch
		if res.err != nil {
			err = res.err
		} else if res.transactions != nil {
			transections = res.transactions
		} else {
			total = res.total
		}
	}
	if err != nil {
		return nil, err
	}
	total_page := (total - offset) / limit

	return &TransactionHistoryModel{
		Transactions:     transections,
		TotalTransection: total,
		TotalPage:        total_page,
	}, nil
}

func (s *Service) publishPaymentEvent(ctx context.Context, tx TransactionModel) error {
	payload, err := json.Marshal(tx)
	if err != nil {
		return nil
	}

	return s.publisher.PublishWithContext(
		ctx,
		"transaction_done",
		"payment.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}
