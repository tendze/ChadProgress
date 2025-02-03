package auth_client

import "errors"

var (
	ErrClientUnavailable = errors.New("auth client is not available")
	ErrUserNotFound      = errors.New("user not found")
)

// UserAuthRequestInterface interface summarizes requests for authorization
type UserAuthRequestInterface interface {
	GetLogin() string
	GetPassword() string
}

// UserRegistrationRequest particular structure for registration requests to auth services
type UserRegistrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserRegistrationResponse particular structure for registration responses to auth services
type UserRegistrationResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	Error  string `json:"error"`
}

// UserLoginRequest particular structure for login requests to auth services
type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserLoginResponse particular structure for login responses to auth services
type UserLoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type UserValidateTokenResponse struct {
	Status    string `json:"status"`
	UserLogin string `json:"user-login"`
	Error     string `json:"error"`
}
