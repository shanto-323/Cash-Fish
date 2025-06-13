package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	worker  sarama.Consumer
	service *Service
}

func NewConsumer(url []string, service *Service) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	worker, err := sarama.NewConsumer(url, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		worker:  worker,
		service: service,
	}, nil
}

func (c *Consumer) NewUserSettlementAccount(topic string) error {
	consumer, err := c.worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	defer consumer.Close()

	for {
		select {
		case msg := <-consumer.Messages():
			{
				fmt.Println(string(msg.Value))
				var uid string
				if err := json.Unmarshal(msg.Value, &uid); err != nil {
					log.Println(err)
					continue
				}
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := c.service.NewEntity(ctx, uid); err != nil {
					log.Println(err)
					continue
				}
				log.Println("New User Created!")
			}
		case err := <-consumer.Errors():
			log.Println(err)
			continue
		}
	}
}
