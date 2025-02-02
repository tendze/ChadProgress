package userservice

import (
	"ChadProgress/internal/auth_client"
	"ChadProgress/internal/models"
	"ChadProgress/storage"
	"context"
	"fmt"
	"log/slog"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
)

type AuthServiceClient interface {
	RegisterUser(ctx context.Context, authReq auth_client.UserAuthRequestInterface) (*auth_client.UserRegistrationResponse, error)
}

type UserService struct {
	storage    storage.Storage
	authClient AuthServiceClient
	log        *slog.Logger
}

func NewUserService(
	storage storage.Storage,
	authServiceClient AuthServiceClient,
	log *slog.Logger,
) *UserService {
	return &UserService{
		storage:    storage,
		authClient: authServiceClient,
		log:        log,
	}
}

// RegisterUser This function returns token from side authorization service and error
func (u *UserService) RegisterUser(email, password, name, role string) (string, error) {
	const op = "services.user.user_service.RegisterUser"
	log := u.log.With(
		slog.String("op", op),
	)

	user, _ := u.storage.GetUser(email)
	if user != nil {
		log.Info("user already exists")
		return "", fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
	}

	newUser := &models.User{
		Email: email,
		Name:  name,
		Role:  role,
	}

	regReq := models.UserAuth{
		Login:    email,
		Password: password,
	}

	resp, err := u.authClient.RegisterUser(context.Background(), regReq)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	userID, err := u.storage.SaveUser(newUser)
	if err != nil {
		if err == storage.ErrUserAlreadyExists {
			log.Info("user already exists")
			return "", ErrUserAlreadyExists
		}
		log.Error("save user failed", slog.String("errorType", err.Error()))
		return "", err
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
		return "", err
	}

	jwtToken := resp.Token
	return jwtToken, nil
}
