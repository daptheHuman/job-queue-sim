package api

import (
	"encoding/json"
	"net/http"

	"github.com/daptheHuman/job-orchestrator-go/internal/rabbitmq"
	"github.com/daptheHuman/job-orchestrator-go/internal/storage"
	"github.com/google/uuid"
)

type API struct {
	Queue   *rabbitmq.RabbitMQ
	Storage *storage.JobRepository
}

func NewAPI(conn *rabbitmq.RabbitMQ, repo *storage.JobRepository) *API {
	return &API{
		Queue:   conn,
		Storage: repo,
	}
}

func (a *API) handleJobSubmission(w http.ResponseWriter, r *http.Request) {
	job := storage.Job{
		ID:     uuid.New(),
		Name:   r.FormValue("name"),
		Status: "Pending",
	}

	if _, err := a.Storage.CreateJob(job); err != nil {
		http.Error(w, "Failed to save job in the database", http.StatusInternalServerError)
		return
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Failed to marshal job", http.StatusInternalServerError)
		return
	}

	if err := a.Queue.Publish(jobData); err != nil {
		http.Error(w, "Failed to publish job to the queue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (a *API) handleJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Path[len("/check-jobs/"):]
	if jobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	status, err := a.Storage.GetJobByID(uuid.MustParse(jobID))
	if err != nil {
		http.Error(w, "Failed to get job in the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
