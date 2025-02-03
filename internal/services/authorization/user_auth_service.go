package userauthservice

import (
	"ChadProgress/internal/auth_client"
	"ChadProgress/internal/models"
	"ChadProgress/internal/services"
	"ChadProgress/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type AuthServiceClient interface {
	RegisterUser(ctx context.Context, authReq auth_client.UserAuthRequestInterface) (*auth_client.UserRegistrationResponse, error)
	LoginUser(ctx context.Context, authReq auth_client.UserAuthRequestInterface) (*auth_client.UserLoginResponse, error)
}

type UserAuthService struct {
	// TODO: MAKE OWN STORAGE INTERFACE
	storage    storage.Storage
	authClient AuthServiceClient
	log        *slog.Logger
}

func NewUserAuthService(
	storage storage.Storage,
	authServiceClient AuthServiceClient,
	log *slog.Logger,
) *UserAuthService {
	return &UserAuthService{
		storage:    storage,
		authClient: authServiceClient,
		log:        log,
	}
}

// RegisterUser This function returns token from side authorization service and error
func (u *UserAuthService) RegisterUser(email, password, name, role string) (string, error) {
	const op = "services.user.user_service.RegisterUser"
	log := u.log.With(
		slog.String("op", op),
	)

	_, err := u.storage.GetUser(email)
	if err == nil {
		log.Info("user already exists")
		return "", fmt.Errorf("%s: %w", op, service.ErrUserAlreadyExists)
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

	_, err = u.storage.SaveUser(newUser)
	if err != nil {
		if errors.Is(err, storage.ErrFieldIsTooLong) {
			log.Info("field is too long")
			return "", fmt.Errorf("%s: %w", op, service.ErrFieldIsTooLong)
		}
		log.Error("save user failed", slog.String("errorType", err.Error()))
		return "", err
	}

	jwtToken := resp.Token
	return jwtToken, nil
}

func (u *UserAuthService) Login(email, password string) (string, error) {
	const op = "services.user.user_service.Login"
	log := u.log.With(
		slog.String("op", op),
	)

	regReq := models.UserAuth{
		Login:    email,
		Password: password,
	}

	loginResp, err := u.authClient.LoginUser(context.Background(), regReq)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user successfully signed in")
	return loginResp.Token, nil
}
