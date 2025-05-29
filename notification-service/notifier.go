package notificationservice

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Notifier(d amqp.Delivery) {
	fmt.Println(d.Body)
	// can use any push notification service
}
