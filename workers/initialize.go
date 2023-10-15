package workers

import (
	"backend-auth/configs/dev"
	"backend-auth/internal/logger"
	"backend-auth/internal/utils/connection"
	"fmt"
	"os"
	"strings"
)

func InitializeWorker(workerFn func()) {
	if os.Getenv("ENV") == "dev" {
		dev.LoadGlobalEnvFile()
		dev.LoadWorkersEnvFile()
	}

	fmt.Printf("Starting worker %s", os.Getenv("WORKER_NAME"))
	var forever chan struct{}

	queuesNames := os.Getenv("Q")
	queues := strings.Split(queuesNames, "-q")
	fmt.Printf("Starting for queues %s", queues)
	worker := connection.InitializeWorker()
	for _, queueName := range queues {
		msgs, err := worker.InitializeConsumer(queueName)
		if err != nil {
			logger.LogFailure(err, fmt.Sprintf("Error initializing consumer for queue %s", queueName))
			panic(err)
		}

		queueName := queueName
		go func() {
			for d := range msgs {
				fmt.Printf("Received a message for queue %s: %s", queueName, d.Body)
				workerFn()
			}
		}()

	}
	<-forever

}
