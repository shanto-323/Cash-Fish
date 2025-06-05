package authservice

import (
	"context"
	"fmt"

	pkg "auth-service/pkg"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUp(ctx context.Context, username, password, email string) (*UserResponseModel, error)
	SignIn(ctx context.Context, email, password string) (*UserResponseModel, error)
	SignOut(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user UserModel) error
	AddCard(ctx context.Context, uid, number, brand string, exp_m, exp_y int) (*[]CardsResponseMetadata, error)
	DeleteUser(ctx context.Context, id string) error
	DeleteCard(ctx context.Context, uid string) error
}

type authService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &authService{
		repo: repo,
	}
}

func (s *authService) SignUp(ctx context.Context, username, password, email string) (*UserResponseModel, error) {
	user := UserModel{
		Username: username,
		Email:    email,
	}
	user.ID = ksuid.New().String()
	h_pass, err := pkg.NewHashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = h_pass
	token, r_token, err := pkg.NewToken(user.ID)
	if err != nil {
		return nil, err
	}
	user.RefreshToken = r_token

	if err = s.repo.NewUser(ctx, user); err != nil {
		return nil, err
	}

	return &UserResponseModel{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Cards:    []CardsResponseMetadata{},
		Token: TokenMetadata{
			Token:        token,
			RefreshToken: r_token,
		},
	}, nil
}

func (s *authService) SignIn(ctx context.Context, email, password string) (*UserResponseModel, error) {
	resp, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if resp.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	if err := pkg.CompareHash(resp.Password, password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("error mismatch password")
		}
		return nil, err
	}

	token, r_token, err := pkg.NewToken(resp.ID)
	if err != nil {
		return nil, err
	}

	s.UpdateUser(ctx, UserModel{
		ID:           resp.ID,
		Username:     resp.Username,
		Password:     resp.Password,
		Email:        resp.Email,
		RefreshToken: r_token,
	})

	resp.Token.Token = token
	resp.Token.RefreshToken = r_token
	return resp, nil
}

func (s *authService) SignOut(ctx context.Context, id string) error {
	return s.repo.UpdateToken(ctx, id, "")
}

func (s *authService) UpdateUser(ctx context.Context, user UserModel) error {
	resp, err := s.repo.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	if mutationHelper(user.Username) {
		resp.Username = user.Username
	}
	if mutationHelper(user.Password) {
		if err := pkg.CompareHash(resp.Password, user.Password); err == bcrypt.ErrMismatchedHashAndPassword {
			hash, err := pkg.NewHashPassword(user.Password)
			if err != nil {
				return err
			}
			resp.Password = hash
		}
	}
	if mutationHelper(user.Email) {
		resp.Email = user.Email
	}
	return s.repo.UpdateUser(ctx, UserModel{
		ID:           resp.ID,
		Username:     resp.Username,
		Password:     resp.Password,
		Email:        resp.Email,
		RefreshToken: resp.Token.RefreshToken,
	})
}

func (s *authService) AddCard(ctx context.Context, uid, number, brand string, exp_m, exp_y int) (*[]CardsResponseMetadata, error) {
	return s.repo.NewCard(ctx, CardMetadata{
		UID:         uid,
		Number:      number,
		Brand:       brand,
		ExpiryMonth: exp_m,
		ExpiryYear:  exp_y,
	})
}

func (s *authService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *authService) DeleteCard(ctx context.Context, uid string) error {
	return s.repo.DeleteCard(ctx, uid)
}

func mutationHelper(value string) bool {
	if value != "" {
		return true
	}
	return false
}
