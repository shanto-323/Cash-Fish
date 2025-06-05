package auth

import (
	"context"
	"fmt"
	"net"

	authservice "auth-service/internal"
	"auth-service/pb"

	"google.golang.org/grpc"
)

type grpcAuthServer struct {
	service authservice.Service
	pb.UnimplementedAuthServiceServer
}

func NewGrpcServer(s authservice.Service, p int) error {
	port := fmt.Sprintf(":%d", p)
	ls, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterAuthServiceServer(serv, &grpcAuthServer{service: s})
	return serv.Serve(ls)
}

func (g *grpcAuthServer) UserSignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	resp, err := g.service.SignUp(ctx, r.User.Username, r.User.Password, r.User.Email)
	if err != nil {
		return nil, err
	}
	return &pb.SignUpResponse{
		Id: resp.ID,
		User: &pb.User{
			Username: resp.Username,
			Password: resp.Password,
			Email:    resp.Email,
		},
		Token:        resp.Token.Token,
		RefreshToken: resp.Token.RefreshToken,
	}, nil
}

func (g *grpcAuthServer) UserSignIn(ctx context.Context, r *pb.SignInRequest) (*pb.SignInResponse, error) {
	resp, err := g.service.SignIn(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}
	return &pb.SignInResponse{
		Id: resp.ID,
		User: &pb.User{
			Username: resp.Username,
			Password: resp.Password,
			Email:    resp.Email,
		},
		Token:        resp.Token.Token,
		RefreshToken: resp.Token.RefreshToken,
	}, nil
}

func (g *grpcAuthServer) UserSignOut(ctx context.Context, r *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	err := g.service.SignOut(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.SignOutResponse{
		Msg: fmt.Sprintf("signed out success:%s\n", r.Id),
	}, nil
}

func (g *grpcAuthServer) UserUpdate(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := g.service.UpdateUser(ctx, authservice.UserModel{
		ID:       r.Id,
		Username: r.User.Username,
		Password: r.User.Password,
		Email:    r.User.Email,
	}); err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		Msg: fmt.Sprintf("user updated id:%s\n", r.Id),
	}, nil
}

func (g *grpcAuthServer) UserDelete(ctx context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := g.service.DeleteUser(ctx, r.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{
		Msg: fmt.Sprintf("user deleted by id:%s\n", r.Id),
	}, nil
}

func (g *grpcAuthServer) UserNewToken(ctx context.Context, r *pb.NewTokenRequest) (*pb.NewTokenResponse, error) {
	resp, err := g.service.NewToken(ctx, r.Id, r.Token)
	if err != nil {
		return nil, err
	}

	return &pb.NewTokenResponse{
		Token: resp,
	}, nil
}
