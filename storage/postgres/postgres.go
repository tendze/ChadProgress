package postgres

import (
	"ChadProgress/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Storage struct {
	DB *gorm.DB
}

func New(dsn string) (*Storage, error) {
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
	return &Storage{DB: db}, nil
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

func (s *Storage) SaveUser(user *models.User) (int64, error) {
	//TODO: обработка ошибки, когда некорректный role
	result := s.DB.Create(user)
	if result.Error != nil {
		log.Fatalf("ошибка при добавлении юзера")
	}
	return int64(user.ID), nil
}

func (s *Storage) SaveClient(client *models.Client) error {
	result := s.DB.Create(client)
	if result.Error != nil {
		log.Fatalf("failed to save trainer")
	}
	return nil
}

func (s *Storage) SaveTrainer(trainer *models.Trainer) error {
	result := s.DB.Create(trainer)
	if result.Error != nil {
		log.Fatalf("failed to save trainer")
	}
	return nil
}

func (s *Storage) GetUser(email string) (*models.Client, error) { return nil, nil }
