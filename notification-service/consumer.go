package notificationservice

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConsumer(rebbitmqUrl string) (*Consumer, error) {
	conn, err := amqp.Dial(rebbitmqUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *Consumer) Consume() {
	defer c.conn.Close()
	defer c.ch.Close()

	queue, err := c.ch.QueueDeclare(
		"tr.status",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Print(err)
	}

	msg, err := c.ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Print(err)
	}

	for d := range msg {
		Notifier(d)
	}
}
