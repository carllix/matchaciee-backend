package database

import (
	"fmt"
	"log"
	"time"

	"github.com/carllix/matchaciee-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error

	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		switch cfg.LogLevel {
		case "debug":
			logLevel = logger.Info
		case "info":
			logLevel = logger.Warn
		default:
			logLevel = logger.Error
		}
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	dsn := cfg.GetDSN()
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func IsConnected() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		return false
	}

	return true
}
