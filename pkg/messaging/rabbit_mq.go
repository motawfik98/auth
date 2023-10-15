package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
)

type RabbitMq struct {
	client *amqp.Connection
}

func initializeRabbitMqConnection() (*amqp.Connection, error) {
	connString := os.ExpandEnv("amqp://${RABBITMQ_USERNAME}:${RABBITMQ_PASSWORD}@${RABBITMQ_HOST}:${RABBITMQ_PORT}/")
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (r *RabbitMq) createQueue(name string) error {
	ch, err := r.client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	fmt.Println(q)
	return err
}

func (r *RabbitMq) SendMessage(queueName string, body map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := r.client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        bodyBytes,
	})
}

func (r *RabbitMq) CreateConsumer(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, err := r.client.Channel()
	if err != nil {
		return nil, nil, err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, nil, err
	}
	return ch, msgs, nil
}
