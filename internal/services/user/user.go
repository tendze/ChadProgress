package userservice

import (
	"ChadProgress/internal/models"
	"errors"
	"fmt"
	"log/slog"
)

type Storage interface {
	GetUser(email string) (*models.User, error)
	SaveTrainer(trainer *models.Trainer) error
}

type UserService struct {
	storage Storage
	log     *slog.Logger
}

func NewUserService(
	storage Storage,
	log *slog.Logger,
) *UserService {
	return &UserService{
		storage: storage,
		log:     log,
	}
}

func (u *UserService) CreateTrainer(userEmail, qualification, experience, achievement string) error {
	const op = "services.user.user.CreateTrainer"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUser(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))
		return errors.New("user not found")
	}

	newTrainer := &models.Trainer{
		UserID:         user.ID,
		Qualifications: qualification,
		Experience:     experience,
		Achievements:   achievement,
	}

	err := u.storage.SaveTrainer(newTrainer)
	if err != nil {
		return err
	}

	return nil
}
