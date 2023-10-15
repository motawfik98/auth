package main

import (
	"backend-auth/cmd/workers"
	workerModel "backend-auth/pkg/workers"
	"fmt"
)

func workerFn(worker *workerModel.Worker) {
	fmt.Println("Received a job")
}

func main() {
	workers.StartWorker(workerFn)
}
