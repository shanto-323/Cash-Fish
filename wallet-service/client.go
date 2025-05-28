package walletservice

import (
	"context"

	"cash-fish/wallet-service/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.WalletServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	service := pb.NewWalletServiceClient(conn)
	return &Client{
		conn:    conn,
		service: service,
	}, nil
}

func (c *Client) CreatePayment(ctx context.Context, senderId, receiverId, note, idempotencyKey string, amount float64) (*TransactionResponseModel, error) {
	resp, err := c.service.CreatePayment(
		ctx,
		&pb.CreatePaymentRequest{
			SenderId:       senderId,
			ReceiverId:     receiverId,
			Note:           note,
			IdempotencyKey: idempotencyKey,
			Amount:         amount,
		},
	)
	if err != nil {
		return nil, err
	}

	return &TransactionResponseModel{
		PaymentId:     resp.PaymentId,
		PaymentStatus: pb.PaymentStatus_PAYMENT_STATUS_COMPLETED.String(),
		CreatedAt:     resp.CreatedAt.AsTime(),
	}, nil
}

func (c *Client) GetPaymentStatus(ctx context.Context, payment_id string) (*TransactionModel, error) {
	resp, err := c.service.GetPaymentStatus(ctx, &pb.GetPaymentStatusRequest{PaymentId: payment_id})
	if err != nil {
		return nil, err
	}

	r := resp.Transection
	id, err := uuid.Parse(r.PaymentId)
	if err != nil {
		return nil, err
	}
	transaction := &TransactionModel{
		ID:         id,
		SenderId:   r.SenderId,
		ReceiverId: r.ReceiverId,
		Note:       r.Note,
		Amount:     int64(r.Amount / 100),
		CreatedAt:  r.CreatedAt.AsTime(),
	}
	return transaction, nil
}

func (c *Client) GetTransectionHistory(ctx context.Context, id string, limit, offset int64) (*TransactionHistoryModel, error) {
	resp, err := c.service.GetTransectionHistory(
		ctx,
		&pb.GetTransectionHistoryRequest{
			UserId: id,
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		return nil, err
	}

	transactions := []*TransactionModel{}
	for _, r := range resp.Transection {
		id, err := uuid.Parse(r.PaymentId)
		if err != nil {
			return nil, err
		}
		transaction := &TransactionModel{
			ID:         id,
			SenderId:   r.SenderId,
			ReceiverId: r.ReceiverId,
			Note:       r.Note,
			Amount:     int64(r.Amount / 100),
			CreatedAt:  r.CreatedAt.AsTime(),
		}
		transactions = append(transactions, transaction)
	}

	return &TransactionHistoryModel{
		Transactions:     transactions,
		TotalTransection: resp.TotalTransection,
		TotalPage:        resp.TotalPage,
	}, nil
}
