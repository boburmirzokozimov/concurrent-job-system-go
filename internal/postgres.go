package config

import (
	"concurrent-job-system/models"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(cfg *Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	gormConfig := &gorm.Config{}
	if cfg.Environment != "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("❌ Failed to get raw SQL DB: %v", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	if err := Migrate(db); err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return nil, err
	}

	log.Println("✅ Connected to Postgres and ran migrations.")
	return &PostgresDB{
		db,
	}, nil
}

func Migrate(p *gorm.DB) error {
	return p.AutoMigrate(
		&models.BaseJob{},
	)
}

func (s *PostgresDB) Save(job models.Processable) error {
	payload, _ := json.Marshal(job) // serialize full struct
	entity := models.BaseJob{
		JobType:       job.Type(),
		Payload:       string(payload),
		Priority:      int(job.GetPriority()),
		MaxRetryCount: job.GetMaxRetryCount(),
		Retries:       job.GetRetries(),
		Status:        "pending",
		CreatedAt:     job.Base().CreatedAt,
		UpdatedAt:     job.Base().UpdatedAt,
	}
	res := s.db.Create(&entity)
	return res.Error
}

func (s *PostgresDB) UpdateStatus(id int, status string) error {
	return s.db.Model(&models.BaseJob{}).Where("id = ?", id).Update("status", status).Error
}

func (s *PostgresDB) LoadPending() ([]models.BaseJob, error) {
	var entities []models.BaseJob
	s.db.Where("status != ?", "success").Find(&entities)

	var jobs []models.BaseJob
	for _, e := range entities {
		var base models.BaseJob
		json.Unmarshal([]byte(e.Payload), &base) // deserialize base
		base.ID = e.ID
		base.Status = e.Status
		jobs = append(jobs, base)
	}
	return jobs, nil
}
