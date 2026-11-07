package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func New(logFilePath string) (*slog.Logger, io.Closer, error) {
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	writer := io.MultiWriter(file, os.Stdout)
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{}))

	return logger, file, nil
}
