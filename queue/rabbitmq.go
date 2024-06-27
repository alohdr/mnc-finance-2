// queue/rabbitmq.go
package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"mnc-finance/config"
	"os"
)

type PublishDefinition interface {
	ProduceMessage(queueName, routingKey string, message []byte) error
}
type PublishService struct {
	rabbit config.RabbitMQ
}

func NewPublishService(rabbit config.RabbitMQ) PublishDefinition {
	return PublishService{
		rabbit: rabbit,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func (p PublishService) ProduceMessage(queueName, routingKey string, message []byte) error {
	exchangeName := os.Getenv("RabbitExchange")
	err := p.rabbit.Channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	_, err = p.rabbit.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange with the specified routing key
	err = p.rabbit.Channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
	err = p.rabbit.Channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	//}
	failOnError(err, "Failed to publish a message")
	fmt.Println("Message published: %s", string(message))
	return nil
}
