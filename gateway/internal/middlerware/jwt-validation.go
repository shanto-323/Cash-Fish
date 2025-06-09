package middlerware

import (
	"fmt"
	"net/http"
	"time"

	"gateway/pkg"

	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
)

type JwtClaims struct {
	ID string
	*jwt.StandardClaims
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := r.Cookie("access_token")
			if err != nil {
				pkg.WriteJson(w, http.StatusBadGateway, fmt.Errorf("token err%s", err))
				return
			}
			if accessToken.Value == "" {
				pkg.WriteJson(w, http.StatusBadRequest, fmt.Errorf("access token is nil%s", err))
				return
			}

			claims, err := validateJwt(accessToken.Value)
			if err != nil {
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorExpired != 0 {
						pkg.WriteJson(w, http.StatusBadGateway, fmt.Errorf("validation error%s", err))
						return
					}
				}
				pkg.WriteJson(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
				return
			}

			id := r.URL.Query().Get("id")
			if id != claims.ID {
				pkg.WriteJson(w, http.StatusBadGateway, fmt.Errorf("id not matching"))
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

type config struct {
	JwtKey string `envconfig:"JWT_KEY"`
}

func validateJwt(token string) (*JwtClaims, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	signedToken, err := jwt.ParseWithClaims(
		token,
		&JwtClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("token encryption method not matching")
			}
			return []byte(cfg.JwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := signedToken.Claims.(*JwtClaims)
	if !ok {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
