package controller

import (
	"gateway/internal/service/auth"

	"github.com/gorilla/mux"
)

type AuthController struct {
	router     *mux.Router
	authClient *auth.AuthClient
}

func NewAuthController(router *mux.Router, authClient *auth.AuthClient) *AuthController {
	router.PathPrefix("/auth").Subrouter()
	return &AuthController{
		router:     router,
		authClient: authClient,
	}
}

func (c *AuthController) RegisterRoutes() {
}
