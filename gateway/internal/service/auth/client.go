package auth

import (
	"context"

	pb "gateway/internal/service/pb"

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

func (a *AuthClient) AuthClientSignUP(ctx context.Context, username, email, password string) (*UserResponseModel, error) {
	resp, err := a.service.SignUp(ctx, &pb.SignUpRequest{
		User: &pb.User{
			Username: username,
			Password: password,
			Email:    email,
		},
	})
	if err != nil {
		return nil, err
	}

	return &UserResponseModel{
		ID:       resp.Id,
		Username: resp.User.Username,
		Password: resp.User.Password,
		Email:    resp.User.Email,
		Cards:    []CardsResponseMetadata{},
		Token: TokenMetadata{
			Token:        resp.Token,
			RefreshToken: resp.RefreshToken,
		},
	}, nil
}

func (a *AuthClient) AuthClientSignIN(ctx context.Context, email, password string) (*UserResponseModel, error) {
	resp, err := a.service.SignIn(ctx, &pb.SignInRequest{
		Password: password,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}

	return &UserResponseModel{
		ID:       resp.Id,
		Username: resp.User.Username,
		Password: resp.User.Password,
		Email:    resp.User.Email,
		Cards:    []CardsResponseMetadata{}, // TOO MUCH PAYLOAD SOLUTION -> []SLICE{EMPTY}
		Token: TokenMetadata{
			Token:        resp.Token,
			RefreshToken: resp.RefreshToken,
		},
	}, nil
}

func (a *AuthClient) AuthClientSignOut(ctx context.Context, id string) (*string, error) {
	resp, err := a.service.SignOut(ctx, &pb.SignOutRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &resp.Msg, nil
}

// FOR NOW USER WILL NOT UPDATE

func (a *AuthClient) AuthClientDeleteUser(ctx context.Context, id string) (*string, error) {
	resp, err := a.service.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &resp.Msg, nil
}

func (a *AuthClient) AuthClientNewToken(ctx context.Context, id, r_token string) (*string, error) {
	resp, err := a.service.NewToken(ctx, &pb.NewTokenRequest{Id: id, Token: r_token})
	if err != nil {
		return nil, err
	}

	return &resp.Token, nil
}
