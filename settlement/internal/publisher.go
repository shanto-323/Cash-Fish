package internal

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TRANSECTION_QUEUE        = "tr.status"
	SETTLEMENT_SERVICE_QUEUE = "se.status"
)

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(rabbitmqUrl string) (*Publisher, error) {
	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *Publisher) CloseConn() error {
	return p.conn.Close()
}

func (p *Publisher) CloseChan() error {
	return p.conn.Close()
}

func (p *Publisher) Consumer(ctx context.Context, sourceQueue string, targetQueue string) error {
	_, err := p.ch.QueueDeclare(
		sourceQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := p.ch.Consume(
		sourceQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := p.Produce(ctx, targetQueue, msg); err != nil {
				log.Print(err)
				continue // keep listening event
			}
		}
	}()
	return nil
}

func (p *Publisher) Produce(ctx context.Context, queueName string, msg amqp.Delivery) error {
	_, err := p.ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
		},
	)
}
