package userservice

import (
	"errors"
	"fmt"
	"log/slog"

	"ChadProgress/internal/models"
	service "ChadProgress/internal/services"
	"ChadProgress/storage"
)

type Storage interface {
	GetUserByEmail(email string) (*models.User, error)
	GetTrainerByID(id uint) (*models.Trainer, error)
	GetTrainerByUserID(id uint) (*models.Trainer, error)
	GetClientByUserID(id uint) (*models.Client, error)
	SaveTrainer(trainer *models.Trainer) error
	SaveClient(client *models.Client) error
	UpdateTrainerID(clientID, trainerID uint) error
	GetTrainersClients(trainerID uint) ([]models.Client, error)
	CreatePlan(plan *models.TrainingPlan) error
	AddMetrics(metric *models.Metric) error
	GetMetrics(clientID uint) ([]models.Metric, error)
	AddProgressReport(report *models.ProgressReport) error
	GetProgressReport(trainerID, clientID uint) ([]models.ProgressReport, error)
	GetPlan(trainerID, clientId uint) ([]models.TrainingPlan, error)
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
		log.Error(fmt.Sprintf("user with email <%s> tried to create trainer profile while being client", userEmail))
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
		log.Error("trainer cant create client profile")

		return service.ErrInvalidRoleRequest
	}

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
		log.Error("trainer cant select trainer")

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
		log.Error("trainer cant get info about client")
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
		log.Error("trainer profile not found")
		return nil, service.ErrTrainerNotFound
	}

	return trainer, nil
}

func (u *UserService) GetTrainersClients(userEmail string) ([]models.Client, error) {
	const op = "services.user.user.GetTrainersClients"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", userEmail))

		return nil, service.ErrUserNotFound
	}
	if user.Role != models.RoleTrainer {
		log.Error("clients do not have clients")

		return nil, service.ErrInvalidRoleRequest
	}

	trainer, _ := u.storage.GetTrainerByUserID(user.ID)
	if trainer == nil {
		log.Error("trainer profile not found")

		return nil, service.ErrTrainerNotFound
	}

	clients, err := u.storage.GetTrainersClients(trainer.ID)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (u *UserService) CreatePlan(trainerEmail string, clientID uint, description, schedule string) error {
	const op = "services.user.user.CreatePlan"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(trainerEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", trainerEmail))

		return service.ErrUserNotFound
	}

	if user.Role != models.RoleTrainer {
		log.Error("clients do not have clients")

		return service.ErrInvalidRoleRequest
	}

	trainer, _ := u.storage.GetTrainerByUserID(user.ID)
	if trainer == nil {
		log.Error("trainer profile not found")

		return service.ErrTrainerNotFound
	}

	plan := models.TrainingPlan{
		TrainerID:   trainer.ID,
		ClientID:    clientID,
		Description: description,
		Schedule:    schedule,
	}

	err := u.storage.CreatePlan(&plan)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) AddMetrics(clientEmail string, height, weight, bodyFat, bmi float64, measuredAt models.CustomTime) error {
	const op = "services.user.user.AddMetrics"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(clientEmail)

	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", clientEmail))
		return service.ErrUserNotFound
	}

	if user.Role != models.RoleClient {
		log.Error("trainers cant add metrics")
		return service.ErrInvalidRoleRequest
	}

	client, err := u.storage.GetClientByUserID(user.ID)
	if err != nil {
		log.Error("client profile not found")

		return service.ErrUserNotFound
	}

	metric := &models.Metric{
		ClientID:   client.ID,
		Height:     height,
		Weight:     weight,
		BodyFat:    bodyFat,
		BMI:        bmi,
		MeasuredAt: measuredAt.Time,
	}

	err = u.storage.AddMetrics(metric)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetMetrics(clientEmail string) ([]models.Metric, error) {
	const op = "services.user.user.AddMetrics"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(clientEmail)

	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", clientEmail))

		return []models.Metric{}, service.ErrUserNotFound
	}

	if user.Role != models.RoleClient {
		log.Error("trainers cant add metrics")

		return []models.Metric{}, service.ErrInvalidRoleRequest
	}

	client, err := u.storage.GetClientByUserID(user.ID)
	if err != nil {
		log.Error("client profile not found")

		return []models.Metric{}, service.ErrUserNotFound
	}

	metrics, err := u.storage.GetMetrics(client.ID)
	if err != nil {
		// TODO: return more detailed error
		return []models.Metric{}, err
	}

	return metrics, nil
}

func (u *UserService) AddProgressReport(trainerEmail, comments string, clientID uint) error {
	const op = "services.user.user.AddProgressReport"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(trainerEmail)

	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", user.Email))

		return service.ErrUserNotFound
	}

	if user.Role != models.RoleTrainer {
		log.Error("client cant add progress report")

		return service.ErrInvalidRoleRequest
	}

	trainer, err := u.storage.GetTrainerByUserID(user.ID)
	if err != nil {
		log.Error("trainer profile not found")
		return service.ErrInvalidRoleRequest
	}

	report := &models.ProgressReport{
		TrainerID: trainer.ID,
		ClientID:  clientID,
		Comments:  comments,
	}

	err = u.storage.AddProgressReport(report)
	if err != nil {
		log.Error("error occurred while adding progress report")

		return err
	}

	return nil
}

func (u *UserService) GetProgressReport(userEmail string, trainerID, clientID uint) ([]models.ProgressReport, error) {
	const op = "services.user.user.AddProgressReport"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", user.Email))

		return []models.ProgressReport{}, service.ErrUserNotFound
	}

	reports, err := u.storage.GetProgressReport(trainerID, clientID)
	if err != nil {
		log.Error("error occurred while getting progress report", slog.String("error", err.Error()))

		return []models.ProgressReport{}, err
	}

	return reports, nil
}

func (u *UserService) GetPlan(userEmail string, trainerID, clientID uint) ([]models.TrainingPlan, error) {
	const op = "services.user.user.GetPlan"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUserByEmail(userEmail)
	if user == nil {
		log.Error(fmt.Sprintf("user with email <%s> not found", user.Email))

		return []models.TrainingPlan{}, service.ErrUserNotFound
	}

	plans, err := u.storage.GetPlan(trainerID, clientID)
	if err != nil {
		log.Error("error occurred while getting plan", slog.String("error", err.Error()))

		return []models.TrainingPlan{}, err
	}

	return plans, nil
}
