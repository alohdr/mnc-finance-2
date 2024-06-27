// queue/rabbitmq.go
package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel
var QueueName = "transferQueue"

func init() {
	var err error
	conn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	_, err = ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
}

func Publish(body []byte) error {
	return ch.Publish(
		"",
		QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func Consume() (<-chan amqp.Delivery, error) {
	return ch.Consume(
		QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func Close() {
	ch.Close()
	conn.Close()
}
