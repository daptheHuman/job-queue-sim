package worker

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/daptheHuman/job-orchestrator-go/internal/rabbitmq"
	"github.com/daptheHuman/job-orchestrator-go/internal/storage"
	"github.com/streadway/amqp"
)

type Worker struct {
	Storage  *storage.JobRepository
	RabbitMQ *rabbitmq.RabbitMQ
	maxJobs  int
}

func NewWorker(storage *storage.JobRepository, rabbitmq *rabbitmq.RabbitMQ) *Worker {
	return &Worker{
		Storage:  storage,
		RabbitMQ: rabbitmq,
		maxJobs:  2,
	}
}

func (w *Worker) StartWorker(n int) {
	for i := 0; i < n; i++ {
		go w.RunWorker(i)
	}

	select {}
}

func (w *Worker) RunWorker(id int) {
	log.Printf("Worker %d started", id)

	msgs, err := w.RabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Worker %d Failed to register consumer: %v", id, err)
	}

	curJobs := make(chan struct{}, w.maxJobs)

	go func() {
		for msg := range msgs {
			if len(curJobs) >= w.maxJobs {
				log.Printf("Worker %d cannot accept more jobs", id)
				msg.Nack(false, true)
				continue
			}

			curJobs <- struct{}{}

			go func(m amqp.Delivery) {
				defer func() {
					<-curJobs
				}()
				job := storage.Job{}
				if err := json.Unmarshal(m.Body, &job); err != nil {
					log.Printf("Failed to unmarshal job: %v", err)
				}

				jobIdString := job.ID.String()
				log.Printf("Worker %d processing job: %s", id, jobIdString[:3])
				w.Storage.UpdateJobStatus(job.ID, "In-Progress")
				if completed := processJob(jobIdString); !completed {
					log.Printf("Worker %d failed to process job: %s, requeuing...", id, jobIdString[:3])
					m.Nack(false, true)
					return
				}
				log.Printf("Worker %d successfully prcoess job: %s", id, jobIdString[:3])
				w.Storage.UpdateJobStatus(job.ID, "Completed")
				m.Ack(false)
			}(msg)
		}
	}()

	select {}
}

func processJob(jobID string) bool {
	time.Sleep(time.Second * time.Duration(3)) // Simulate a delay to represent job processing
	// Chance to simulate job failure
	if rand.Intn(10)%2 == 0 {
		log.Printf("Failed to process job: %s", jobID[:3])
		return false
	}

	return true
}
