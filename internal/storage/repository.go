package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) CreateJob(payload Job) (uuid.UUID, error) {
	query := `INSERT INTO jobs (id, name, status) 
              VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, payload.ID, payload.Name, payload.Status)
	if err != nil {
		log.Fatalf("Failed to insert job into the database: %v", err)
		return payload.ID, err
	}
	return payload.ID, nil
}

func (r *JobRepository) UpdateJobStatus(id uuid.UUID, status string) error {
	query := `UPDATE jobs SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		log.Fatalf("Failed to update job status in the database: %v", err)
	}
	return err
}

func (r *JobRepository) GetJobByID(id uuid.UUID) (*Job, error) {
	job := &Job{}
	query := `SELECT id, status, created_at, updated_at FROM jobs WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&job.ID, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return job, nil
}
