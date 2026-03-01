package sqlx

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSQLDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &config.Config)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to get sql db: %w", err)
	}

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 100
	}

	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 1 * time.Hour
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	fmt.Printf("[postgres] connected to postgres: host=%s, port=%s, user=%s, db=%s\n", config.Host, config.Port, config.User, config.DBName)

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("[postgres] failed to get sql db: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("[postgres] failed to close sql db: %w", err)
	}

	return nil
}
