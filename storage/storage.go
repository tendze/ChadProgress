package storage

import "ChadProgress/internal/models"

type Storage interface {
	SaveUser(user *models.User) (int64, error)
	SaveClient(client *models.Client) error
	SaveTrainer(trainer *models.Trainer) error
	GetUser(email string) (*models.Client, error)
	// All methods
}
