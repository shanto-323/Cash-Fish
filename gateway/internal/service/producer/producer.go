package producer

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(url []string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 2

	producer, err := sarama.NewSyncProducer(url, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) PushToQueue(message any, topic string) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	newMessage := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	_, _, err = p.producer.SendMessage(&newMessage)
	if err != nil {
		return err
	}

	return nil
}
