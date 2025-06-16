package db

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/factory"
	"concurrent-job-system/internal/logger"
	"encoding/json"
)

type JobRepository struct {
	db     *PostgresDB
	logger logger.ILogger
}

func NewJobRepository(db *PostgresDB, logger logger.ILogger) *JobRepository {
	return &JobRepository{
		db:     db,
		logger: logger,
	}
}

type IJobRepository interface {
	Save(j job.IProcessable) (int, error)
	UpdateStatus(id int, status string) error
	LoadPending() (job.IProcessable, error)
}

func (r *JobRepository) Save(j job.IProcessable) (int, error) {
	payload, _ := json.Marshal(j.PayloadOnly())
	entity := job.BaseJob{
		JobType:       j.Type(),
		Payload:       string(payload),
		Priority:      int(j.GetPriority()),
		MaxRetryCount: j.GetMaxRetryCount(),
		Retries:       j.GetRetries(),
		Status:        "pending",
		CreatedAt:     j.Base().CreatedAt,
		UpdatedAt:     j.Base().UpdatedAt,
	}
	res := r.db.db.Create(&entity)
	if res.Error != nil {
		r.logger.Error("Failed to save job %v: %v", j.GetId(), res.Error)
	}
	return entity.ID, res.Error
}

func (r *JobRepository) UpdateStatus(id int, status string) error {
	err := r.db.db.Model(&job.BaseJob{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		r.logger.Error("Failed to update job %d status to %s: %v", id, status, err)
	}
	return err
}

func (r *JobRepository) LoadPending() (job.IProcessable, error) {
	var row job.BaseJob
	if err := r.db.db.Where("status != ?", "success").Limit(1).Find(&row).Error; err != nil {
		return nil, err
	}

	j, err := factory.DeserializeJob(row)
	if err != nil {
		r.logger.Warn("skipping corrupt job id %d: %v", row.ID, err)
	}
	return j, nil
}
