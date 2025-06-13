package config

import (
	"concurrent-job-system/models"
	"encoding/json"
	"time"
)

type IStorage interface {
	Save(j models.Processable) error
	UpdateStatus(id int, status string) error
	LoadPending() ([]models.BaseJob, error)
}
type JobRecord struct {
	ID            int
	Type          string
	Priority      string
	Retries       int
	MaxRetryCount int
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Payload       json.RawMessage
}
