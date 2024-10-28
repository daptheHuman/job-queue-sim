package rabbitmq

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func Connect() (*RabbitMQ, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"jobs",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil

}

func (r *RabbitMQ) Publish(message []byte) error {
	return r.Channel.Publish(
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(
		r.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
