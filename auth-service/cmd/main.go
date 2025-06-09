package main

import (
	"log"
	"time"

	authservice "auth-service/internal"
	auth "auth-service/internal/auth"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseDsn       string `envconfig:"DATABASE_DSN"`
	AuthServiceIpAddr string `envconfig:"AUTH_SERVICE_IP_ADDR"`
	CardServiceIpAddr string `envconfig:"CARD_SERVICE_IP_ADDR"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var repository authservice.Repository
	var err error
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			repository, err = authservice.NewRepository(cfg.DatabaseDsn)
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		},
	)
	service := authservice.NewService(repository)
	log.Println("auth server running on port", cfg.AuthServiceIpAddr)
	log.Fatal(auth.NewGrpcServer(service, cfg.AuthServiceIpAddr))
}
