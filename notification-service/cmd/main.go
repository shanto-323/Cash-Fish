package main

import (
	"log"
	"time"

	co "cash-fish/notification-service"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	RabbitmqUrl string `envconfig:"RABBITMQ_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("could not get env variabls", err)
	}

	var consumer *co.Consumer
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			consumer, err = co.NewConsumer(cfg.RabbitmqUrl)
			if err != nil {
				return err
			}
			return nil
		},
	)

	forever := make(chan interface{})
	go func() {
		consumer.Consume()
	}()
	<-forever
}
