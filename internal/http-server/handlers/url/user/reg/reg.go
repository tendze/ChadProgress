package reg

import (
	"ChadProgress/internal/lib/api/response"
	"ChadProgress/internal/models"
	userservice "ChadProgress/internal/services/user"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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

	var req Request
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", err.Error())
		render.JSON(w, r, response.Error("failed to decode request body"))
		return
	}
	log.Info("request body decoded", slog.Any("request", req))
	if err = validator.New().Struct(req); err != nil {
		log.Error("invalid request")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request"))
		return
	}

	//TODO: pass a hashed password!!!
	err = u.userService.RegisterUser(req.Email, req.Password, req.Name, req.Role)

	if err != nil {
		log.Error("failed to save user", err.Error())
		w.WriteHeader(http.StatusBadGateway)
		render.JSON(w, r, response.Error("failed to save user"))
		return
	}
	log.Info("successfully saved user", slog.String("email", req.Email))
	render.JSON(w, r, response.OK())
}
