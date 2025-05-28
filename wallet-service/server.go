package walletservice

import (
	"context"
	"fmt"
	"net"

	"cash-fish/wallet-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcServer struct {
	pb.UnsafeWalletServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterWalletServiceServer(serv, &grpcServer{service: s})
	return serv.Serve(ln)
}

func (sr *grpcServer) CreatePayment(ctx context.Context, r *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	resp, err := sr.service.CreateNewTransaction(
		ctx,
		r.SenderId,
		r.ReceiverId,
		r.Note,
		r.IdempotencyKey,
		r.Amount,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId: resp.ID.String(),
		Status:    pb.PaymentStatus_PAYMENT_STATUS_COMPLETED,
		CreatedAt: timestamppb.New(resp.CreatedAt),
	}, nil
}

func (sr *grpcServer) GetPaymentStatus(ctx context.Context, r *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {
	resp, err := sr.service.GetTransaction(ctx, r.PaymentId)
	if err != nil {
		return nil, err
	}

	transection := &pb.Transection{
		PaymentId:  resp.ID.String(),
		SenderId:   resp.SenderId,
		ReceiverId: resp.ReceiverId,
		Note:       resp.Note,
		Amount:     float64(resp.Amount / 100),
		CreatedAt:  timestamppb.New(resp.CreatedAt),
	}

	return &pb.GetPaymentStatusResponse{
		Transection: transection,
	}, nil
}

func (sr *grpcServer) GetTransectionHistory(ctx context.Context, r *pb.GetTransectionHistoryRequest) (*pb.GetTransectionHistoryResponse, error) {
	resp, err := sr.service.GetTransactionHistory(ctx, r.UserId, r.Limit, r.Offset)
	if err != nil {
		return nil, err
	}

	transections := []*pb.Transection{}
	for _, r := range resp.Transactions {
		transection := &pb.Transection{
			PaymentId:  r.ID.String(),
			SenderId:   r.SenderId,
			ReceiverId: r.ReceiverId,
			Note:       r.Note,
			Amount:     float64(r.Amount / 100),
			CreatedAt:  timestamppb.New(r.CreatedAt),
		}
		transections = append(transections, transection)
	}

	return &pb.GetTransectionHistoryResponse{
		Transection:      transections,
		TotalTransection: resp.TotalTransection,
		TotalPage:        resp.TotalPage,
	}, nil
}
