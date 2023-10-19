package workers

import (
	"backend-auth/configs/dev"
	"backend-auth/internal/logger"
	"backend-auth/internal/utils/connection"
	"backend-auth/pkg/workers"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"strings"
)

func StartWorker(workerFn func(worker *workers.Worker, delivery amqp.Delivery) error) {
	if os.Getenv("ENV") == "dev" {
		dev.LoadGlobalEnvFile()
		dev.LoadWorkersEnvFile()
	}

	fmt.Printf("Starting worker %s\n", os.Getenv("WORKER_NAME"))
	var forever chan struct{}

	queuesNames := os.Getenv("Q")
	queues := strings.Split(queuesNames, "-q ")
	fmt.Printf("Starting for queues %s\n", queues)
	worker := connection.InitializeWorker()
	for _, queueName := range queues {
		if len(queueName) == 0 {
			continue
		}
		ch, msgs, err := worker.InitializeConsumer(queueName)
		if err != nil {
			logger.LogFailure(err, fmt.Sprintf("Error initializing consumer for queue %s", queueName))
			panic(err)
		}

		go func(queueName string, ch *amqp.Channel, msgs <-chan amqp.Delivery) {
			for d := range msgs {
				fmt.Printf("Received a message for queue %s: %s\n", queueName, d.Body)
				if err := workerFn(worker, d); err == nil {
					d.Ack(true)
				} else {
					fmt.Printf("====>>> ERROR FOR MESSAGE: %s", d.Body)
					d.Ack(true)
				}
			}
			defer ch.Close()
		}(queueName, ch, msgs)

	}
	<-forever

}
