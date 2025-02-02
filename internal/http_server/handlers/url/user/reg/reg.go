package reg

import (
	"ChadProgress/internal/lib/api/response"
	"ChadProgress/internal/models"
	userservice "ChadProgress/internal/services/user"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=trainer client"`
}

type RegisterResponse struct {
	Status   string      `json:"status"`
	JWTToken string      `json:"token"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type UserProvider interface {
	SaveUser(email, password, name, role string) error
	GetUser(email string) *models.Client
}

type UserHandler struct {
	userService *userservice.UserService
	log         *slog.Logger
}

func NewUserHandler(
	// TODO: INTERFACE INSTEAD OF STRUCT
	service *userservice.UserService,
	log *slog.Logger,
) *UserHandler {
	return &UserHandler{userService: service, log: log}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
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
		log.Error("failed to decode request body", err.Error())
		render.JSON(w, r, response.Error("failed to decode request body"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))
	if err = validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		validationErr := err.(validator.ValidationErrors)
		log.Error("invalid request", validationErr.Error())
		render.JSON(w, r, response.ValidationError(validationErr))
		return
	}

	jwtToken, err := u.userService.RegisterUser(req.Email, req.Password, req.Name, req.Role)

	if err != nil {
		if errors.Is(err, userservice.ErrUserAlreadyExists) {
			log.Info("user already exists")
			w.WriteHeader(http.StatusBadGateway)
			render.JSON(w, r, response.Error("user already with such email"))
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

func regResponseOK(token string) RegisterResponse {
	return RegisterResponse{
		Status:   "OK",
		JWTToken: token,
	}
}
