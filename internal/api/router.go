package api

import (
	"net/http"

	"github.com/daptheHuman/job-orchestrator-go/internal/rabbitmq"
	"github.com/daptheHuman/job-orchestrator-go/internal/storage"
)

func SetupRouter(rabbitmq *rabbitmq.RabbitMQ, storage *storage.JobRepository) *http.ServeMux {
	mux := http.NewServeMux()
	api := NewAPI(rabbitmq, storage)
	mux.HandleFunc("/jobs", api.handleJobSubmission)
	mux.HandleFunc("/check-jobs/", api.handleJobStatus)

	return mux
}
