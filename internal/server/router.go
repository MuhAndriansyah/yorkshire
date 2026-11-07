package server

import (
	"net/http"

	"github.com/MuhAndriansyah/yorkshire/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Pasang semua routes dari Handler ke Chi Router
	handler.Mount(r, h)

	return r
}
