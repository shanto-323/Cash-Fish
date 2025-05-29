package walletservice

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PAYMENT_STATUS_COMPLETE = "COMPLETE"
	PAYMENT_STATUS_ERROR    = "ERROR"
	PAYMENT_STATUS_PENDING  = "PENDING"
	PAYMENT_STATUS_REVERSED = "REVERSED"
)

type ProducerModel struct {
	PaymentId  string
	SenderId   string
	ReceiverId string
	Status     string
}

type Producer struct {
	publisher *amqp.Channel
}

func NewProducer(ch *amqp.Channel) *Producer {
	return &Producer{publisher: ch}
}

func (pr *Producer) Produce(ctx context.Context, p ProducerModel) error {
	queue, err := pr.publisher.QueueDeclare(
		"tr.status",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return pr.publisher.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}
