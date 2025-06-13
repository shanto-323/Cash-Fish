package main

import (
	"log"
	"time"

	"settlement/internal"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseDsn     string   `envconfig:"DATABASE_DSN"`
	BrokerUrl       []string `envconfig:"BROKER_URL"`
	UserEntityTopic string   `envconfig:"USER_ENTITY_TOPIC"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var (
		repository internal.Repository
		consumer   *internal.Consumer
		err        error
	)

	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			repository, err = internal.NewRepository(cfg.DatabaseDsn)
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		},
	)

	service := internal.NewService(repository)

	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			consumer, err = internal.NewConsumer(cfg.BrokerUrl, service)
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		},
	)

	log.Println("Settlement Service running")

	forever := make(chan interface{})
	go func() {
		consumer.NewUserSettlementAccount(cfg.UserEntityTopic)
	}()
	<-forever
}
