package workers

import (
	"backend-auth/internal/logger"
	controllerUtil "backend-auth/utils/controller"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func InitializeWorker(workerFn func()) {
	if os.Getenv("ENV") == "dev" {
		if err := godotenv.Load("../..", "."); err != nil {
			panic(err)
		}
	}

	fmt.Printf("Starting worker %s", os.Getenv("WORKER_NAME"))
	var forever chan struct{}

	queuesNames := os.Getenv("Q")
	queues := strings.Split(queuesNames, "-q")
	fmt.Printf("Starting for queues %s", queues)
	controller := controllerUtil.InitializeController()
	for _, queueName := range queues {
		msgs, err := controller.InitializeConsumer(queueName)
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
