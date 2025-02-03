package userhandler

import (
	"ChadProgress/internal/lib/api/response"
	"ChadProgress/internal/models"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type CreateTrainerProfileRequest struct {
	Qualification string `json:"qualification" validate:"required"`
	Experience    string `json:"experience" validate:"required"`
	Achievement   string `json:"achievement" validate:"required"`
}

type CreateTrainerProfileResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type UserService interface {
	CreateTrainer(userEmail, qualification, experience, achievement string) error
}

type UserHandler struct {
	log         *slog.Logger
	userService UserService
}

func NewUserHandler(
	log *slog.Logger,
	userService UserService,
) *UserHandler {
	return &UserHandler{log: log, userService: userService}
}

func (u *UserHandler) CreateTrainer(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.CreateTrainer"
	log := u.log.With(
		slog.String("op", op),
	)
	userEmail := r.Context().Value(models.ContextUserKey).(string)
	if userEmail == "" {
		log.Error("empty email from context")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	var req CreateTrainerProfileRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", err.Error())
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("could not decode request body"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))
	if err = validator.New().Struct(req); err != nil {
		validationErr := err.(validator.ValidationErrors)
		log.Error("invalid request", validationErr.Error())
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.ValidationError(validationErr))
		return
	}

	log.Info("user email extracted from context", slog.String("email", userEmail))
	err = u.userService.CreateTrainer(userEmail, req.Qualification, req.Experience, req.Achievement)

	if err != nil {
		log.Error("create trainer failed")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("create trainer failed"))
		return
	}

	setHeaderRenderJSON(w, r, http.StatusOK, response.OK())
}

func setHeaderRenderJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.WriteHeader(status)
	render.JSON(w, r, v)
}
