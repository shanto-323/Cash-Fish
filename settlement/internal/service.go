package internal

import (
	"context"
	"log"

	"settlement/pkg"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) NewEntity(ctx context.Context, uid string) error {
	newEntity := pkg.UserBalance{
		UserId:  uid,
		Balance: 0,
		Active:  true,
	}
	err := s.repo.MakeEntry(ctx, newEntity)
	if err != nil {
		return err
	}

	log.Println("USER ENTITY CREATED")
	return nil
}

func (s *Service) TransfarMoney(ctx context.Context) error {
	return nil
}

func (s *Service) CashIn(ctx context.Context) error {
	return nil
}

func (s *Service) CashOut(ctx context.Context) error {
	return nil
}
