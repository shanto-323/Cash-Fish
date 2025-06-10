package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gateway/internal/controller"
	"gateway/internal/service/auth"
	"gateway/internal/service/card"
	p "gateway/internal/service/producer"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type config struct {
	BrokerUrl         []string `envconfig:"BROKER_URL"`
	UserCreationTopic string   `envconfig:"USER_CREATION_TOPIC"`
}

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
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}
	var producer *p.Producer
	var err error
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			producer, err = p.NewProducer(cfg.BrokerUrl, cfg.UserCreationTopic)
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		},
	)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	authRouter := controller.NewAuthController(apiRouter, s.authClient, producer)
	authRouter.RegisterRoutes()

	cardRouter := controller.NewCardController(apiRouter, s.cardClient)
	cardRouter.RegisterRoutes()

	fmt.Println("api running on port", s.ipAddr)
	ip := fmt.Sprintf(":%s", s.ipAddr)
	return http.ListenAndServe(ip, router)
}
