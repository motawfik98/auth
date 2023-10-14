package controllers

import amqp "github.com/rabbitmq/amqp091-go"

func (c *Controller) InitializeConsumer(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := c.messaging.Connection.CreateConsumer(queueName)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
