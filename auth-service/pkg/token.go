package pkg

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
)

type JwtClaims struct {
	ID string
	*jwt.StandardClaims
}

type Config struct {
	KEY string `envconfig:"JWT_KEY"`
}

func NewToken(id string) (string, string, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return "", "", err
	}

	claim := &JwtClaims{
		ID: id,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Hour).Unix(),
		},
	}
	t, err := token(claim)
	if err != nil {
		return "", "", err
	}

	r_claim := &JwtClaims{
		ID: id,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Hour).Unix(),
		},
	}
	r_token, err := token(r_claim)
	if err != nil {
		return "", "", err
	}

	return t, r_token, nil
}

func token(claim *JwtClaims) (string, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(cfg.KEY))
}
