package settlementservice

import "context"

type Service struct {
	repo      Repository
	publisher *Publisher
}

func NewService(repo Repository, publisher *Publisher) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) NewEntity(ctx context.Context, b UserBalance) error {
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
