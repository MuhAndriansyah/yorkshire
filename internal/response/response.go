package response

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/MuhAndriansyah/yorkshire/internal/apperror"
)

type Success struct {
	Status bool `json:"status"`
	Data   any  `json:"data"`
	Meta   any  `json:"meta,omitempty"`
}

type Error struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type CursorPagination struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	Limit      int    `json:"limit"`
	HasMore    bool   `json:"has_more"`
}

type OffsetPagination struct {
	Page       int `json:"page"`
	TotalPage  int `json:"total_page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_item"`
}

func writeJson(w http.ResponseWriter, statusCode int, val any) {
	w.Header().Set("Content-Application", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(val); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func JSON(w http.ResponseWriter, statusCode int, data, meta any) {
	writeJson(w, statusCode, Success{
		Status: true,
		Data:   data,
		Meta:   meta,
	})
}

func ErrorMsg(w http.ResponseWriter, statusCode int, message string) {
	writeJson(w, statusCode, Error{
		Status:  false,
		Message: message,
	})
}

func ValidationError(w http.ResponseWriter, statusCode int, message string, errs map[string]string) {
	writeJson(w, statusCode, Error{
		Status:  false,
		Message: message,
		Errors:  errs,
	})
}

func HandleError(w http.ResponseWriter, logger *slog.Logger, err error) {
	if appErr, ok := errors.AsType[*apperror.AppError](err); ok {
		if appErr.Code >= http.StatusInternalServerError {
			logger.Error(appErr.Message, "error", err)
		} else {
			logger.Warn(appErr.Message, "error", err)
		}

		ErrorMsg(w, appErr.Code, appErr.Message)
		return
	}

	if errors.Is(err, context.Canceled) {
		logger.Info("request canceled by client")
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		logger.Warn("request deadline exceeded", "error", err)
		ErrorMsg(w, http.StatusGatewayTimeout, "request timed out")
		return
	}

	logger.Error("unhandled internal error", "error", err)
	ErrorMsg(w, http.StatusInternalServerError, "internal server error")

}
