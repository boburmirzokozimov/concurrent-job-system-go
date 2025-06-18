package db

import (
	"concurrent-job-system/internal/config"
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/logger"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db     *gorm.DB
	logger logger.ILogger
}

func NewPostgresDB(cfg *config.Config, log logger.ILogger) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Error("Failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to get raw SQL DB: %v", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	if err := Migrate(db); err != nil {
		log.Error("Migration failed: %v", err)
		return nil, err
	}

	log.Info("Connected to Postgres and ran migrations.")
	return &PostgresDB{db: db, logger: log}, nil
}

func Migrate(p *gorm.DB) error {
	return p.AutoMigrate(
		&job.BaseJob{},
	)
}
