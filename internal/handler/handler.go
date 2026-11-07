package handler

import (
	"log/slog"

	"github.com/MuhAndriansyah/yorkshire/internal/api"
	"github.com/MuhAndriansyah/yorkshire/internal/service/auth"
	"github.com/go-chi/chi/v5"
	govalidator "github.com/go-playground/validator/v10"
)

const basePath = "/api/v1"

// Handler menggabungkan seluruh sub-handler domain (Auth, User, dll.)
type Handler struct {
	*AuthHandler
}

func New(logger *slog.Logger, validate *govalidator.Validate, authSvc auth.AuthService) *Handler {
	return &Handler{
		AuthHandler: NewAuthHandler(logger, validate, authSvc),
	}
}

// Mount mendaftarkan semua routes yang digenerate oleh oapi-codegen ke Chi Router.
func Mount(r chi.Router, h *Handler) {
	api.HandlerFromMuxWithBaseURL(h, r, basePath)
}

// Memastikan Handler memenuhi api.ServerInterface
var _ api.ServerInterface = (*Handler)(nil)
