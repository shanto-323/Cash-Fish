package main

import (
	"gateway/internal/server"
	"gateway/internal/service/auth"
	"gateway/internal/service/card"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type config struct {
	GatewayIp      string `envconfig:"GATEWAY_IP"`
	AuthServiceUrl string `envconfig:"AUTH_SERVICE_URL"`
	CardServiceUrl string `envconfig:"CARD_SERVICE_URL"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	var (
		err        error
		authClient *auth.AuthClient
		cardClient *card.CardClient
	)

	retry.ForeverSleep(
		2*time.Second,
		func(i int) error {
			authClient, err = auth.NewAuthClient(cfg.AuthServiceUrl)
			cardClient, err = card.NewCardhClient(cfg.CardServiceUrl)
			if err != nil {
				return err
			}
			return nil
		},
	)

	server, err := server.NewServer(cfg.GatewayIp, authClient, cardClient)
	if err != nil {
		log.Panic(err)
	}
	server.StartServer()
}
