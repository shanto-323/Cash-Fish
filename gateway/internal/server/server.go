package server

import (
	"fmt"
	"net/http"

	"gateway/internal/controller"
	"gateway/internal/service/auth"
	"gateway/internal/service/card"

	"github.com/gorilla/mux"
)

type Server struct {
	ipAddr     string
	authClient *auth.AuthClient
	cardClient *card.CardClient
}

func NewServer(authUrl, cardUrl string) (*Server, error) {
	return &Server{}, nil
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()
	router.PathPrefix("/api/v1")

	authRouter := controller.NewAuthController(router, s.authClient)
	authRouter.RegisterRoutes()

	cardRouter := controller.NewCardController(router, s.cardClient)
	cardRouter.RegisterRoutes()

	ip := fmt.Sprintf(":%s", s.ipAddr)
	return http.ListenAndServe(ip, router)
}
