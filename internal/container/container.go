package container

import (
	"concurrent-job-system/internal/config"
	"concurrent-job-system/internal/db"
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/types"
	"concurrent-job-system/internal/logger"
	"concurrent-job-system/internal/worker"
)

type Container struct {
	Cfg           *config.Config
	Pool          *worker.Pool
	Logger        logger.ILogger
	JobRepository db.IJobRepository
}

func (c *Container) MakeJob(i int) job.IProcessable {
	if i%2 == 0 {
		return types.NewSimpleJob()
	} else {
		return types.NewExcelJob("excel.path")
	}
}

func NewContainer() (*Container, error) {
	cfg := config.NewConfig()
	log := logger.NewConsoleLogger()

	dbInstance, err := db.NewPostgresDB(cfg, log)
	if err != nil {
		return nil, err
	}

	jobRepo := db.NewJobRepository(dbInstance, log)

	pool := worker.NewPool(5, jobRepo, log)

	return &Container{
		Cfg:           cfg,
		Logger:        log,
		Pool:          pool,
		JobRepository: jobRepo,
	}, nil
}
