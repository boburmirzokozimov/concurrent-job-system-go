package db

import (
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/factory"
	"concurrent-job-system/internal/logger"
	"concurrent-job-system/pkg/dto"
	"database/sql"
	"encoding/json"
	"time"
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
	UpdateFailedStatus(id int, status string, errorMessage string) error
	LoadPending() (job.IProcessable, error)
	Load(id int) (job.IProcessable, error)
	GetStats() (dto.JobStats, error)
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

func (r *JobRepository) UpdateFailedStatus(id int, status string, errorMessage string) error {
	failedAt := time.Now()
	err := r.db.db.Model(&job.BaseJob{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"error":     errorMessage,
		"failed_at": &failedAt,
	}).Error

	if err != nil {
		r.logger.Error("Failed to update job %d failed status to %s: %v", id, status, err)
	}
	return err
}

func (r *JobRepository) Load(id int) (job.IProcessable, error) {
	var row job.BaseJob
	err := r.db.db.Model(&job.BaseJob{}).Where("id = ?", id).First(&row).Error
	if err != nil {
		r.logger.Error("Failed to load job %d: %v", id, err)
		return nil, err
	}

	j, err := factory.DeserializeJob(row)
	if err != nil {
		r.logger.Warn("Corrupt job id %d: %v", id, err)
		return nil, err
	}

	return j, nil
}

func (r *JobRepository) LoadPending() (job.IProcessable, error) {
	tx := r.db.db.Begin()
	if tx.Error != nil {
		r.logger.Error("Failed to begin transaction: %v", tx.Error)
		return nil, tx.Error
	}

	var row job.BaseJob
	err := tx.Raw(`
        UPDATE base_jobs
        SET status = 'reserved',
            updated_at = NOW()
        WHERE id = (
            SELECT id
            FROM base_jobs
            WHERE status = 'pending'
            ORDER BY priority DESC, created_at
            FOR UPDATE SKIP LOCKED
            LIMIT 1
        )
        RETURNING *
    `).Scan(&row).Error

	if err != nil {
		r.logger.Error("Failed to reserve job: %v", err)
		_ = tx.Rollback()
		return nil, err
	}

	if row.ID == 0 {
		r.logger.Debug("No pending jobs found to reserve")
		_ = tx.Rollback()
		return nil, nil
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Failed to commit transaction for job ID %d: %v", row.ID, err)
		return nil, err
	}

	j, err := factory.DeserializeJob(row)
	if err != nil {
		r.logger.Warn("Deserialization failed for job ID %d: %v", row.ID, err)
		return nil, err
	}

	r.logger.Info("Reserved job ID %d of type '%s' with priority %d", row.ID, row.JobType, row.Priority)
	return j, nil
}

func (r *JobRepository) GetStats() (dto.JobStats, error) {
	var stats dto.JobStats

	rows, err := r.db.db.Model(&job.BaseJob{}).Select(
		"COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending",
		"COUNT(CASE WHEN status = 'running' THEN 1 END) as running",
		"COUNT(CASE WHEN status = 'succeeded' THEN 1 END) as succeeded",
		"COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed",
	).Rows()

	if err != nil {
		return stats, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("Sth went wrong", err)
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&stats.Pending, &stats.Running, &stats.Succeeded, &stats.Failed); err != nil {
			return stats, err
		}
	}

	return stats, nil
}
