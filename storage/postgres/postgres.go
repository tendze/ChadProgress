package postgres

import (
	"ChadProgress/internal/models"
	"ChadProgress/storage"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	err = createRoleEnum(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = createTrainerStatusEnum(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = autoMigrate(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = createDummyTrainer(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{DB: db}, nil
}

func createRoleEnum(db *gorm.DB) error {
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

func createTrainerStatusEnum(db *gorm.DB) error {
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM pg_type 
			WHERE typname = 'status'
		);
	`).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		return db.Exec(`CREATE TYPE status AS ENUM ('ACTIVE', 'BUSY', 'ON_VACATION');`).Error
	}

	return nil
}

// createDummyTrainer required init function. generates first dummy trainer that every new client will link to by default
func createDummyTrainer(db *gorm.DB) error {
	var dummyTrainer models.Trainer
	if err := db.Where("id = 1").First(&dummyTrainer).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		user := models.User{
			Email: "dummytrainer@mail.ru",
			Name:  "Dummy Trainer",
			Role:  "trainer",
		}
		if err = db.Create(&user).Error; err != nil {
			return err
		}

		defaultTrainer := models.Trainer{
			UserID:         user.ID,
			Qualifications: "I'm dummy!",
			Experience:     "I'm dummy!",
			Achievements:   "I'm dummy!",
			Status:         models.StatusActive,
		}
		if err = db.Create(&defaultTrainer).Error; err != nil {
			return err
		}
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
	const op = "postgres.SaveUser"
	result := s.DB.Create(user)
	if err := result.Error; err != nil {
		if isInvalidEnumError(err) {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrUserAlreadyExists)
		} else if isTooLongFieldError(err) {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrFieldIsTooLong)
		}
		return -1, fmt.Errorf("%s: %w", op, result.Error)
	}
	return int64(user.ID), nil
}

func (s *Storage) SaveClient(client *models.Client) error {
	const op = "postgres.SaveClient"
	result := s.DB.Create(client)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return fmt.Errorf("%s: %w", op, storage.ErrDuplicateKey)
		} else if isTooLongFieldError(result.Error) {
			return fmt.Errorf("%s: %w", op, storage.ErrFieldIsTooLong)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func (s *Storage) SaveTrainer(trainer *models.Trainer) error {
	const op = "postgres.SaveTrainer"
	result := s.DB.Create(trainer)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return fmt.Errorf("%s: %w", op, storage.ErrDuplicateKey)
		} else if isTooLongFieldError(result.Error) {
			return fmt.Errorf("%s: %w", op, storage.ErrFieldIsTooLong)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func (s *Storage) GetUserByEmail(email string) (*models.User, error) {
	const op = "postgres.GetUser"
	var user models.User
	result := s.DB.First(&user, "email = ?", email)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (s *Storage) GetTrainerByID(id uint) (*models.Trainer, error) {
	const op = "postgres.GetTrainerByID"
	var trainer models.Trainer
	result := s.DB.First(&trainer, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &trainer, nil
}

func (s *Storage) GetTrainerByUserID(userID uint) (*models.Trainer, error) {
	const op = "postgres.GetTrainerByUserID"
	var trainer models.Trainer
	result := s.DB.First(&trainer, "user_id = ?", userID)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &trainer, nil
}

func (s *Storage) GetClientByID(id uint) (*models.Client, error) {
	const op = "postgres.GetClientByID"
	var client models.Client
	result := s.DB.First(&client, "id = ?", id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &client, nil
}

func (s *Storage) GetClientByUserID(userID uint) (*models.Client, error) {
	const op = "postgres.GetClientByUserID"
	var client models.Client
	result := s.DB.First(&client, "user_id = ?", userID)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &client, nil
}

func (s *Storage) UpdateTrainerID(clientID, trainerID uint) error {
	// TODO: return more detailed error
	return s.DB.Model(&models.Client{}).Where("id = ?", clientID).Update("trainer_id", trainerID).Error
}

func (s *Storage) GetTrainersClients(trainerID uint) ([]models.Client, error) {
	var clients []models.Client
	res := s.DB.Where("trainer_id = ?", trainerID).Find(&clients)
	if res.Error != nil {
		return []models.Client{}, res.Error
	}
	return clients, nil
}

func (s *Storage) CreatePlan(plan *models.TrainingPlan) error {
	const op = "postgres.CreatePlan"
	result := s.DB.Create(plan)
	if result.Error != nil {
		if isTooLongFieldError(result.Error) {
			return fmt.Errorf("%s: %w", op, storage.ErrFieldIsTooLong)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func (s *Storage) AddMetrics(metric *models.Metric) error {
	const op = "postgres.AddMetrics"
	result := s.DB.Create(metric)
	if err := result.Error; err != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func isInvalidEnumError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 22P02")
}

func isTooLongFieldError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 22001")
}

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}
