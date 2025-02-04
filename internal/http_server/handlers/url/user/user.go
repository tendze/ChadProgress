package userhandler

import (
	"ChadProgress/internal/lib/api/response"
	"ChadProgress/internal/models"
	service "ChadProgress/internal/services"
	"errors"
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

type CreateClientRequest struct {
	Height  float64 `json:"height"`
	Weight  float64 `json:"weight"`
	BodyFat float64 `json:"bodyfat"`
}

type SelectTrainerRequest struct {
	TrainerID uint `json:"trainer-id" validate:"required"`
}

type GetClientProfileResponse struct {
	Height  float64 `json:"height"`
	Weight  float64 `json:"weight"`
	BodyFat float64 `json:"bodyfat"`
}

type GetTrainerProfileResponse struct {
	Qualification string `json:"height"`
	Experience    string `json:"experience"`
	Achievements  string `json:"achievements"`
}

type UserService interface {
	CreateTrainer(userEmail, qualification, experience, achievement string) error
	CreateClient(userEmail string, height, weight, bodyFatPercent float64) error
	SelectTrainer(userEmail string, trainerID uint) error
	GetClientProfile(userEmail string) (*models.Client, error)
	GetTrainerProfile(userEmail string) (*models.Trainer, error)
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
		if errors.Is(err, service.ErrDuplicateKey) {
			log.Error("trainer already exists")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("trainer already exists"))
			return
		} else if errors.Is(err, service.ErrFieldIsTooLong) {
			log.Error("one of fields is too long")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("too long field"))
			return
		} else if errors.Is(err, service.ErrInvalidRoleRequest) {
			log.Error("invalid role request")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("cannot create trainer profile while being client"))
			return
		}
		log.Error("create trainer failed: " + err.Error())
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("create trainer failed"))
		return
	}

	setHeaderRenderJSON(w, r, http.StatusOK, response.OK())
}

func (u *UserHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.CreateClient"
	log := u.log.With(
		slog.String("op", op),
	)
	userEmail := r.Context().Value(models.ContextUserKey).(string)
	if userEmail == "" {
		log.Error("empty email from context")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	var req CreateClientRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", err.Error())
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("could not decode request body"))
		return
	}
	if req.Height < 0.0 || req.Weight < 0.0 || req.BodyFat < 0.0 {
		log.Info("invalid request. negative parameter")
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("negative parameter"))
		return
	}

	log.Info("user email extracted from context", slog.String("email", userEmail))
	err = u.userService.CreateClient(userEmail, req.Height, req.Weight, req.BodyFat)

	if err != nil {
		if errors.Is(err, service.ErrDuplicateKey) {
			log.Error("client already exists")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("trainer already exists"))
			return
		} else if errors.Is(err, service.ErrFieldIsTooLong) {
			log.Error("one of fields is too long")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("too long field"))
			return
		} else if errors.Is(err, service.ErrInvalidRoleRequest) {
			log.Error("invalid role request")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("cannot create client profile while being trainer"))
			return
		}

		log.Error("failed to save client " + err.Error())
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("failed to save client"))
		return
	}

	setHeaderRenderJSON(w, r, http.StatusOK, response.OK())
}

func (u *UserHandler) SelectTrainer(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.user.SelectTrainer"
	log := u.log.With(
		slog.String("op", op),
	)

	userEmail := r.Context().Value(models.ContextUserKey).(string)
	if userEmail == "" {
		log.Error("empty email from context")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	var req SelectTrainerRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", err.Error())
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("could not decode request body"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))
	if err = validator.New().Struct(req); err != nil {
		validationErr := err.(validator.ValidationErrors)
		log.Error("invalid request", slog.String("errormsg", validationErr.Error()))
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.ValidationError(validationErr))
		return
	}

	if req.TrainerID <= 0 {
		log.Error("invalid trainer id to bind")
		setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("invalid trainer id to bind"))
		return
	}

	err = u.userService.SelectTrainer(userEmail, req.TrainerID)
	if err != nil {
		if errors.Is(err, service.ErrClientNotFound) {
			log.Error("client's profile does not exist")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("client's profile does not exist"))
			return
		} else if errors.Is(err, service.ErrTrainerNotFound) {
			log.Error("trainer's profile does not exist")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("trainer's profile does not exist"))
			return
		} else if errors.Is(err, service.ErrNotActiveTrainer) {
			log.Error("not active trainer")
			setHeaderRenderJSON(w, r, http.StatusBadRequest, response.Error("trainer is busy or on vacation"))
			return
		}

		log.Error("failed to bind client to trainer")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("failed to bind client to trainer"))
		return
	}

	setHeaderRenderJSON(w, r, http.StatusOK, response.OK())
}

func (u *UserHandler) GetClientProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.user.GetClientProfile"
	log := u.log.With(
		slog.String("op", op),
	)

	userEmail := r.Context().Value(models.ContextUserKey).(string)
	if userEmail == "" {
		log.Error("empty email from context")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	client, err := u.userService.GetClientProfile(userEmail)
	if err != nil {
		log.Error("failed to get client profile")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	clientResp := GetClientProfileResponse{
		BodyFat: client.BodyFat,
		Height:  client.Height,
		Weight:  client.Weight,
	}
	setHeaderRenderJSON(w, r, http.StatusOK, clientResp)
}

func (u *UserHandler) GetTrainerProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.url.user.user.GetTrainerProfile"
	log := u.log.With(
		slog.String("op", op),
	)

	userEmail := r.Context().Value(models.ContextUserKey).(string)
	if userEmail == "" {
		log.Error("empty email from context")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	trainer, err := u.userService.GetTrainerProfile(userEmail)
	if err != nil {
		log.Error("failed to get trainer profile")
		setHeaderRenderJSON(w, r, http.StatusBadGateway, response.Error("bad gateway"))
		return
	}

	clientResp := GetTrainerProfileResponse{
		Qualification: trainer.Qualifications,
		Experience:    trainer.Experience,
		Achievements:  trainer.Achievements,
	}
	setHeaderRenderJSON(w, r, http.StatusOK, clientResp)
}

func setHeaderRenderJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.WriteHeader(status)
	render.JSON(w, r, v)
}
