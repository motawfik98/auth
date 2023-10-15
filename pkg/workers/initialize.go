package workers

import (
	"backend-auth/pkg/cache"
	"backend-auth/pkg/database"
	"backend-auth/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker struct {
	datasource *database.DB
	cache      *cache.Cache
	messaging  *messaging.Messaging
}

func (w *Worker) SetDatasource(db *database.DB) {
	w.datasource = db
}

func (w *Worker) SetCache(cache *cache.Cache) {
	w.cache = cache
}

func (w *Worker) SetMessaging(messaging *messaging.Messaging) {
	w.messaging = messaging
}

func (w *Worker) InitializeConsumer(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, msgs, err := w.messaging.Connection.CreateConsumer(queueName)
	if err != nil {
		return nil, nil, err
	}
	return ch, msgs, nil
}
