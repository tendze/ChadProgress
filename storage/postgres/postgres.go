package postgres

import (
	"ChadProgress/internal/models"
	"ChadProgress/storage"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
)

type Storage struct {
	DB *gorm.DB
}

func New(dsn string) (*Storage, error) {
	const op = "postgres.New"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Отключаем логирование
	})
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
	result := s.DB.Create(user)
	if err := result.Error; err != nil {
		if isInvalidEnum(err) {
			return -1, storage.ErrUserAlreadyExists
		}
		return -1, result.Error
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

func (s *Storage) GetUser(email string) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, "email = ?", email)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func isInvalidEnum(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}
