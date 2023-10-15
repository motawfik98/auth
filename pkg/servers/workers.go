package servers

import amqp "github.com/rabbitmq/amqp091-go"

func (s *Server) InitializeConsumer(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := s.messaging.Connection.CreateConsumer(queueName)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
