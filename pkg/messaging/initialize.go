package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"strconv"
)

type iMessaging interface {
	createQueue(name string) error
	SendMessage(queueName string, body map[string]interface{}) error
	CreateConsumer(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error)
}

type Messaging struct {
	Connection iMessaging
}

func (m *Messaging) InitializeConnection() error {
	rabbitMqEnabled, _ := strconv.ParseBool(os.Getenv("RABBITMQ_ENABLED"))
	if rabbitMqEnabled {
		rabbitMqConnection, err := initializeRabbitMqConnection()
		if err != nil {
			return err
		}
		rabbitMq := new(RabbitMq)
		rabbitMq.client = rabbitMqConnection
		m.Connection = rabbitMq
	}
	return nil
}

func (m *Messaging) CreateQueues() (string, error) {
	err := m.Connection.createQueue("auth::invalidate-refresh-token-family")
	if err != nil {
		return "auth::invalidate-refresh-token-family", err
	}
	return "", nil
}
