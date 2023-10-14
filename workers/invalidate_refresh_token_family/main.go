package main

import (
	"backend-auth/workers"
	"fmt"
)

func workerFn() {
	fmt.Println("Received a job")
}

func main() {
	workers.InitializeWorker(workerFn)
}
