package storage

import (
	"ChadProgress/internal/models"
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrFieldIsTooLong    = errors.New("field is too long")
	ErrDuplicateKey      = errors.New("duplicate key value violates unique constraint")
)

type Storage interface {
	SaveUser(user *models.User) (int64, error)
	SaveClient(client *models.Client) error
	SaveTrainer(trainer *models.Trainer) error
	GetUser(email string) (*models.User, error)
	// All methods
}
