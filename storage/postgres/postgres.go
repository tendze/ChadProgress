package postgres

import (
	"ChadProgress/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ()

func New(dsn string) (*gorm.DB, error) {
	const op = "postgres.New"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = createEnum(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = autoMigrate(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return db, nil
}

func createEnum(db *gorm.DB) error {
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM pg_type 
			WHERE typname = 'role_enum'
		);
	`).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		return db.Exec(`CREATE TYPE role_enum AS ENUM ('trainer', 'client');`).Error
	}

	return nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Trainer{},
		&models.Client{},
		&models.TrainingPlan{},
		&models.ProgressReport{},
		&models.Metric{},
	)
}
