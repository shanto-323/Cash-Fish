package card

import (
	"context"
	"fmt"
	"net"

	authservice "auth-service/internal"
	"auth-service/pb"

	"google.golang.org/grpc"
)

type grpCardServer struct {
	service authservice.Service
	pb.UnimplementedCardServiceServer
}

func NewGrpcServer(s authservice.Service, p string) error {
	port := fmt.Sprintf(":%s", p)
	ls, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterCardServiceServer(serv, &grpCardServer{service: s})
	return serv.Serve(ls)
}

func (g *grpCardServer) AddCard(ctx context.Context, r *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	resp, err := g.service.AddCard(ctx, r.Uid, r.Card.Number, r.Card.Brand, int(r.Card.ExpMonth), int(r.Card.ExpYear))
	if err != nil {
		return nil, err
	}

	cards := []*pb.Card{}
	for _, c := range *resp {
		card := &pb.Card{
			Id:       c.ID,
			Number:   c.Number,
			Brand:    c.Brand,
			ExpMonth: int32(c.ExpiryMonth),
			ExpYear:  int32(c.ExpiryYear),
		}
		cards = append(cards, card)
	}
	return &pb.AddCardResponse{
		Card: cards, // RESPONSE PAYLOAD CAN BE OVERLOAD
	}, nil
}

func (g *grpCardServer) DeleteCards(ctx context.Context, r *pb.DeleteCardsRequest) (*pb.DeleteCardsResponse, error) {
	if err := g.service.DeleteAllCard(ctx, r.Uid); err != nil {
		return nil, err
	}

	return &pb.DeleteCardsResponse{
		Msg: fmt.Sprintln("deleted all card"),
	}, nil
}

func (g *grpCardServer) GetCards(ctx context.Context, r *pb.GetCardsRequest) (*pb.GetCardsResponse, error) {
	resp, err := g.service.GetAllCard(ctx, r.Uid)
	if err != nil {
		return nil, err
	}

	cards := []*pb.Card{}
	for _, c := range *resp {
		card := &pb.Card{
			Id:       c.ID,
			Number:   c.Number,
			Brand:    c.Brand,
			ExpMonth: int32(c.ExpiryMonth),
			ExpYear:  int32(c.ExpiryYear),
		}
		cards = append(cards, card)
	}

	return &pb.GetCardsResponse{Card: cards}, nil
}

func (g *grpCardServer) RemoveCard(ctx context.Context, r *pb.RemoveCardRequest) (*pb.RemoveCardResponse, error) {
	if err := g.service.RemoveCard(ctx, r.Uid, r.Id); err != nil {
		return nil, err
	}

	return &pb.RemoveCardResponse{
		Msg: fmt.Sprintln("card removed"),
	}, nil
}
