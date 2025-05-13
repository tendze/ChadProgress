package authorization

import (
	"errors"
	"log/slog"
	"net/http"

	"ChadProgress/internal/lib/api/response"
	service "ChadProgress/internal/services"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -source=authorization.go -destination=./authorization_mock.go -package=authorization
type UserAuthService interface {
	RegisterUser(email, password, name, role string) (string, error)
	Login(email, password string) (string, error)
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=trainer client"`
}

type RegisterResponse struct {
	Status   string `json:"status"`
	JWTToken string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Status   string `json:"status"`
	JWTToken string `json:"token,omitempty"`
}

type UserAuthHandler struct {
	userService UserAuthService
	log         *slog.Logger
}

func NewUserAuthHandler(
	service UserAuthService,
	log *slog.Logger,
) *UserAuthHandler {
	return &UserAuthHandler{userService: service, log: log}
}

func (u *UserAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.reg.Register"
	log := u.log.With(
		slog.String("op", op),
	)
	if r.Body == nil || r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("empty request"))
		return
	}

	var req RegisterRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("failed to decode request body"))

		return
	}

	log.Info("register request body decoded", slog.Any("request", req))
	if err = validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		validationErr := err.(validator.ValidationErrors)

		log.Error("invalid request", slog.String("error", validationErr.Error()))
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.ValidationError(validationErr))

		return
	}

	jwtToken, err := u.userService.RegisterUser(req.Email, req.Password, req.Name, req.Role)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			log.Info("user already exists")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("user already with such email"))

			return
		} else if errors.Is(err, service.ErrFieldIsTooLong) {
			log.Info("field login or password is too long")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("login and password must be no more than 100 symbols"))

			return
		}

		log.Error("failed to save user", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadGateway)
		render.JSON(w, r, response.Error("failed to save user"))
		return
	}

	log.Info("successfully saved user", slog.String("email", req.Email))
	render.JSON(w, r, regResponseOK(jwtToken))
}

func (u *UserAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.reg.Login"
	log := u.log.With(
		slog.String("op", op),
	)

	if r.Body == nil || r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("empty request"))
		return
	}

	var req LoginRequest

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("failed to decode request body"))

		return
	}

	if err = validator.New().Struct(req); err != nil {
		validationErr := err.(validator.ValidationErrors)
		log.Error("invalid request", slog.String("error", validationErr.Error()))
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.ValidationError(validationErr))

		return
	}
	log.Info("login request body decoded", slog.Any("request", req))

	jwtToken, err := u.userService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			log.Info("invalid credentials")
			setHeaderRenderJSON(w, r, http.StatusUnauthorized, response.Error("invalid credentials"))
			return
		}
		log.Error("failed to sign in")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("failed to sign in"))

		return
	}

	log.Info("user successfully signed in")
	setHeaderRenderJSON(
		w, r,
		http.StatusOK,
		loginResponseOK(jwtToken),
	)
}

func regResponseOK(token string) RegisterResponse {
	return RegisterResponse{
		Status:   "OK",
		JWTToken: token,
	}
}

func loginResponseOK(token string) LoginResponse {
	return LoginResponse{
		Status:   "OK",
		JWTToken: token,
	}
}

func setHeaderRenderJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.WriteHeader(status)
	render.JSON(w, r, v)
}
