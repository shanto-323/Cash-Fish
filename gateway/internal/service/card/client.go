package card

import (
	"context"

	"gateway/internal/service/card/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CardClient struct {
	conn    *grpc.ClientConn
	service pb.CardServiceClient
}

func NewCardhClient(url string) (*CardClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	service := pb.NewCardServiceClient(conn)
	return &CardClient{
		conn:    conn,
		service: service,
	}, nil
}

func (c *CardClient) CardClientAddCard(ctx context.Context, uid, number, brand string, exp_m, exp_y int32) ([]*CardsResponseMetadata, error) {
	resp, err := c.service.AddCard(ctx, &pb.AddCardRequest{
		Uid: uid,
		Card: &pb.Card{
			Number:   number,
			Brand:    brand,
			ExpMonth: exp_m,
			ExpYear:  exp_y,
		},
	})
	if err != nil {
		return nil, err
	}

	cards := []*CardsResponseMetadata{}
	for _, c := range resp.Card {
		card := &CardsResponseMetadata{
			UID:         c.Id, // THATS CARD_ID
			Number:      c.Number,
			Brand:       c.Brand,
			ExpiryMonth: int(c.ExpMonth),
			ExpiryYear:  int(c.ExpYear),
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (c *CardClient) CardClientGetALlCards(ctx context.Context, uid string) ([]*CardsResponseMetadata, error) {
	resp, err := c.service.GetCards(ctx, &pb.GetCardsRequest{Uid: uid})
	if err != nil {
		return nil, err
	}

	cards := []*CardsResponseMetadata{}
	for _, c := range resp.Card {
		card := &CardsResponseMetadata{
			UID:         c.Id, // THATS CARD_ID
			Number:      c.Number,
			Brand:       c.Brand,
			ExpiryMonth: int(c.ExpMonth),
			ExpiryYear:  int(c.ExpYear),
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (c *CardClient) CardClientRemoveCard(ctx context.Context, uid, id string) (*string, error) {
	resp, err := c.service.RemoveCard(ctx, &pb.RemoveCardRequest{Uid: uid, Id: id}) // ID -> CARD_ID
	if err != nil {
		return nil, err
	}

	return &resp.Msg, nil
}

func (c *CardClient) CardClientDeleteCards(ctx context.Context, uid string) (*string, error) {
	resp, err := c.service.DeleteCards(ctx, &pb.DeleteCardsRequest{Uid: uid})
	if err != nil {
		return nil, err
	}

	return &resp.Msg, nil
}
