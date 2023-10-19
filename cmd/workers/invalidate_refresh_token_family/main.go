package main

import (
	"backend-auth/cmd/workers"
	workerModel "backend-auth/pkg/workers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func workerFn(worker *workerModel.Worker, delivery amqp.Delivery) error {
	return worker.InvalidateCompromisedRefreshTokens(delivery.Body)
}

func main() {
	workers.StartWorker(workerFn)
}
