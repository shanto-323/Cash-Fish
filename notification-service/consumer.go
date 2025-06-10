package notificationservice

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	worker sarama.Consumer
}

func NewConsumer(url []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	consumer, err := sarama.NewConsumer(url, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		worker: consumer,
	}, nil
}

func (c *Consumer) OutputMessage(topic string) error {
	consumer, err := c.worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	defer consumer.Close()

	for {
		select {
		case err := <-consumer.Errors():
			log.Println(err)
			continue
		case msg := <-consumer.Messages():
			Notifier(msg.Value, topic)
		}
	}
}
