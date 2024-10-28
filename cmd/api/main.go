package main

import (
	"fmt"
	"log"
	"net/http"

	api "github.com/daptheHuman/job-orchestrator-go/internal/api"
	"github.com/daptheHuman/job-orchestrator-go/internal/rabbitmq"
	storage "github.com/daptheHuman/job-orchestrator-go/internal/storage"
)

func main() {
	repo, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	rabbitmq, err := setupRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	mux := api.SetupRouter(rabbitmq, repo)
	startServer(mux)
}

func setupDatabase() (*storage.JobRepository, error) {
	dbConn := storage.NewDBConn()
	if err := storage.Migrate(dbConn); err != nil {
		return nil, err
	}
	fmt.Println("Database migrated successfully")

	jobRepo := storage.NewJobRepository(dbConn)
	return jobRepo, nil
}

func setupRabbitMQ() (*rabbitmq.RabbitMQ, error) {
	rabbitmqConn, err := rabbitmq.Connect()
	if err != nil {
		return nil, err
	}
	return rabbitmqConn, nil
}

func startServer(handler http.Handler) {
	log.Println("Starting API server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
