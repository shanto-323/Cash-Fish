package auth

import (
	"auth-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	service pb.AuthServiceClient
}

func NewAuthClient(url string) (*AuthClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	service := pb.NewAuthServiceClient(conn)
	return &AuthClient{
		conn:    conn,
		service: service,
	}, nil
}
