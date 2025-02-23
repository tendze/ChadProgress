package userservice

import (
	"ChadProgress/internal/models"
	service "ChadProgress/internal/services"
	"ChadProgress/storage"
	"errors"
	"fmt"
	"log/slog"
)

type Storage interface {
	GetUserByEmail(email string) (*models.User, error)
	GetTrainerByID(id uint) (*models.Trainer, error)
	GetTrainerByUserID(id uint) (*models.Trainer, error)
	GetClientByUserID(id uint) (*models.Client, error)
	SaveTrainer(trainer *models.Trainer) error
	SaveClient(client *models.Client) error
	UpdateTrainerID(clientID, trainerID uint) error
	GetTrainersClients(trainerID uint) ([]*models.Client, error)
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

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))
		return errors.New("user not found")
	}
	if user.Role == models.RoleClient {
		log.Error(fmt.Sprintf("user with email <%s> tried to create trainer profile while being client"))
		return service.ErrInvalidRoleRequest
	}

	newTrainer := &models.Trainer{
		UserID:         user.ID,
		Qualifications: qualification,
		Experience:     experience,
		Achievements:   achievement,
	}

	err := u.storage.SaveTrainer(newTrainer)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateKey) {
			return service.ErrDuplicateKey
		} else if errors.Is(err, storage.ErrFieldIsTooLong) {
			return fmt.Errorf("%s: %w", op, service.ErrFieldIsTooLong)
		}
		return err
	}

	return nil
}

func (u *UserService) CreateClient(userEmail string, height, weight, bodyFat float64) error {
	const op = "services.user.user.CreateClient"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))
		return errors.New("user not found")
	}
	if user.Role == models.RoleTrainer {
		log.Error(fmt.Sprintf("trainer cant create client profile"))
		return service.ErrInvalidRoleRequest
	}

	//TODO: make one default row in database, чтобы все новые ссылались на него
	newClient := &models.Client{
		UserID:    user.ID,
		TrainerID: 1,
		Height:    height,
		Weight:    weight,
		BodyFat:   bodyFat,
	}

	err := u.storage.SaveClient(newClient)

	if err != nil {
		if errors.Is(err, storage.ErrDuplicateKey) {
			return service.ErrDuplicateKey
		} else if errors.Is(err, storage.ErrFieldIsTooLong) {
			return fmt.Errorf("%s: %w", op, service.ErrFieldIsTooLong)
		}
		return err
	}

	return nil
}

func (u *UserService) SelectTrainer(userEmail string, trainerID uint) error {
	const op = "services.user.user.SelectTrainer"
	log := u.log.With(
		slog.String("op", op),
	)

	clientUser, _ := u.storage.GetUserByEmail(userEmail)
	if clientUser == nil {
		log.Error(fmt.Sprintf("profile with email <%s> not found", userEmail))
		return service.ErrUserNotFound
	}
	if clientUser.Role != models.RoleClient {
		log.Error(fmt.Sprintf("trainer cant select trainer"))
		return service.ErrInvalidRoleRequest
	}

	client, _ := u.storage.GetClientByUserID(clientUser.ID)
	if client == nil {
		log.Error(fmt.Sprintf("client profile with email <%s> not found", userEmail))
		return service.ErrClientNotFound
	}

	trainer, err := u.storage.GetTrainerByID(trainerID)
	if err != nil {
		return service.ErrTrainerNotFound
	}
	if trainer.Status != models.StatusActive {
		return service.ErrNotActiveTrainer
	}

	err = u.storage.UpdateTrainerID(client.ID, trainerID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetClientProfile(userEmail string) (*models.Client, error) {
	const op = "services.user.user.GetClientProfile"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))
		return nil, service.ErrUserNotFound
	}
	if user.Role != models.RoleClient {
		log.Error(fmt.Sprintf("trainer cant get info about client"))
		return nil, service.ErrInvalidRoleRequest
	}

	client, _ := u.storage.GetClientByUserID(user.ID)
	if client == nil {
		return nil, service.ErrClientNotFound
	}

	return client, nil
}

func (u *UserService) GetTrainerProfile(userEmail string) (*models.Trainer, error) {
	const op = "services.user.user.GetTrainerProfile"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))
		return nil, service.ErrUserNotFound
	}

	trainer, _ := u.storage.GetTrainerByUserID(user.ID)
	if trainer == nil {
		log.Error(fmt.Sprintf("trainer profile not found"))
		return nil, service.ErrTrainerNotFound
	}

	return trainer, nil
}

func (u *UserService) GetTrainersClients(userEmail string) ([]*models.Client, error) {
	const op = "services.user.user.GetTrainersClients"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))
		return nil, service.ErrUserNotFound
	}

	trainer, _ := u.storage.GetTrainerByUserID(user.ID)
	if trainer == nil {
		log.Error(fmt.Sprintf("trainer profile not found"))
		return nil, service.ErrTrainerNotFound
	}
	clients, err := u.storage.GetTrainersClients(trainer.ID)
	if err != nil {
		return nil, err
	}
	return clients, nil
}
