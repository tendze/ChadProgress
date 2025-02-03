package storage

import (
	"ChadProgress/internal/models"
	"fmt"
)

var (
	ErrUserAlreadyExists  = fmt.Errorf("user already exists")
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrFieldIsTooLong     = fmt.Errorf("field is too long")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
)

type Storage interface {
	SaveUser(user *models.User) (int64, error)
	SaveClient(client *models.Client) error
	SaveTrainer(trainer *models.Trainer) error
	GetUser(email string) (*models.User, error)
	// All methods
}
