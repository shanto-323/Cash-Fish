package server

import (
	"fmt"
	"net/http"

	"gateway/internal/controller"
	"gateway/internal/service/auth"
	"gateway/internal/service/card"
	"gateway/pkg"

	"github.com/gorilla/mux"
)

type Server struct {
	ipAddr     string
	authClient *auth.AuthClient
	cardClient *card.CardClient
}

func NewServer(ip string, authClient *auth.AuthClient, cardClient *card.CardClient) (*Server, error) {
	return &Server{
		ipAddr:     ip,
		authClient: authClient,
		cardClient: cardClient,
	}, nil
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		pkg.WriteJson(w, http.StatusOK, "hi from server")
	})

	authRouter := controller.NewAuthController(apiRouter, s.authClient)
	authRouter.RegisterRoutes()

	//cardRouter := controller.NewCardController(apiRouter, s.cardClient)
	//cardRouter.RegisterRoutes()

	fmt.Println("api running on port", s.ipAddr)
	ip := fmt.Sprintf(":%s", s.ipAddr)
	return http.ListenAndServe(ip, router)
}
