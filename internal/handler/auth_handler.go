package handler

import (
	"log/slog"
	"net/http"

	"github.com/MuhAndriansyah/yorkshire/internal/api"
	"github.com/MuhAndriansyah/yorkshire/internal/response"
	"github.com/MuhAndriansyah/yorkshire/internal/service/auth"
	"github.com/MuhAndriansyah/yorkshire/internal/validator"
	govalidator "github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	logger       *slog.Logger
	validate     *govalidator.Validate
	authServices auth.AuthService
}

func NewAuthHandler(logger *slog.Logger, validate *govalidator.Validate, authServices auth.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:       logger,
		validate:     validate,
		authServices: authServices,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterJSONRequestBody

	if err := decodeJSONBody(w, r, &req); err != nil {
		response.ErrorMsg(w, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	if err := a.validate.Struct(req); err != nil {
		response.ValidationError(w, http.StatusUnprocessableEntity, "validation failed", validator.ErrorMap(err))
		return
	}

	ctx := r.Context()

	user, err := a.authServices.Register(ctx, auth.RegisterInput{
		Username: req.Username,
		Email:    string(req.Email),
		Password: req.Password,
	})

	if err != nil {
		response.HandleError(w, a.logger, err)
		return
	}

	response.JSON(w, http.StatusCreated, user, nil)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req api.LoginJSONRequestBody

	if err := decodeJSONBody(w, r, &req); err != nil {
		response.ErrorMsg(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := a.validate.Struct(req); err != nil {
		response.ValidationError(w, http.StatusUnprocessableEntity, "validation failed", validator.ErrorMap(err))
		return
	}

	ctx := r.Context()

	user, err := a.authServices.Login(ctx, auth.LoginInput{
		Email:    string(req.Email),
		Password: req.Password,
	})

	if err != nil {
		response.HandleError(w, a.logger, err)
		return
	}

	response.JSON(w, http.StatusOK, user, nil)
}
