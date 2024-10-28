package main

import (
	"log"

	"github.com/daptheHuman/job-orchestrator-go/internal/rabbitmq"
	"github.com/daptheHuman/job-orchestrator-go/internal/storage"
	worker "github.com/daptheHuman/job-orchestrator-go/internal/worker"
)

func main() {
	dbConn := storage.NewDBConn()
	jobRepo := storage.NewJobRepository(dbConn)

	rabbitmqConn, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmqConn.Close()

	workerPool := worker.NewWorker(jobRepo, rabbitmqConn)
	workerPool.StartWorker(2)
}
