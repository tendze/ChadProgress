package userservice

import (
	"ChadProgress/internal/models"
	"ChadProgress/storage"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
)

type UserService struct {
	storage storage.Storage
	log     *slog.Logger
}

func NewUserService(
	storage storage.Storage,
	log *slog.Logger,
) *UserService {
	return &UserService{storage: storage, log: log}
}

func (u *UserService) RegisterUser(email, password, name, role string) error {
	const op = "services.user.user_service.RegisterUser"
	log := u.log.With(
		slog.String("op", op),
	)

	_, err := u.storage.GetUser(email)
	if err == nil {
		if err == storage.ErrUserAlreadyExists {
			log.Info("user already exists")
		} else {
			log.Error("get user failed", slog.String("errorType", err.Error()))
		}
		return errors.New("user already exists")
	}

	newUser := &models.User{
		Email:        email,
		PasswordHash: password,
		Name:         name,
		Role:         role,
	}
	userID, err := u.storage.SaveUser(newUser)
	if err != nil {
		if err == storage.ErrUserAlreadyExists {
			log.Info("user already exists")
			return ErrUserAlreadyExists
		}
		log.Error("save user failed", slog.String("errorType", err.Error()))
		return err
	}

	if role == "client" {
		newClient := &models.Client{
			UserID:    uint(userID),
			TrainerID: 0,
			Height:    0,
			Weight:    0,
			BodyFat:   0,
		}
		err = u.storage.SaveClient(newClient)
	} else if role == "trainer" {
		newTrainer := &models.Trainer{
			UserID:         uint(userID),
			Qualifications: "",
			Experience:     "",
			Achievements:   "",
		}
		err = u.storage.SaveTrainer(newTrainer)
	}

	if err != nil {
		log.Error(fmt.Sprintf("could not save new %s", role))
		return err
	}

	return nil
}
