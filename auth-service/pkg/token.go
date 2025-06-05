package authservice

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	KEY string `envconfig:"JWT_KEY"`
}

func NewToken(id string) (string, string, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return "", "", err
	}

	claim := &jwt.StandardClaims{
		Id:        id,
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}
	t, err := token(claim)
	if err != nil {
		return "", "", err
	}

	r_claim := &jwt.StandardClaims{
		Id:        id,
		ExpiresAt: time.Now().Add(168 * time.Hour).Unix(),
	}
	r_token, err := token(r_claim)
	if err != nil {
		return "", "", err
	}

	return t, r_token, nil
}

func token(claim *jwt.StandardClaims) (string, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(cfg.KEY)
}
