package main

import (
	"log"
	"time"

	notificationservice "notification"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	BrokerUrl         string `envconfig:"BROKER_URL"`
	UserCreationTopic string `envconfig:"USER_CREATION_TOPIC"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("could not get env variabls", err)
	}

	var consumer *notificationservice.Consumer
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			consumer, err = notificationservice.NewConsumer([]string{cfg.BrokerUrl})
			if err != nil {
				return err
			}
			return nil
		},
	)

	forever := make(chan interface{})
	go func() {
		log.Fatal(consumer.OutputMessage(cfg.UserCreationTopic))
	}()
	<-forever
}
