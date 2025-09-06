package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func NewLogger(env string, w io.Writer) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case "dev":
		logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "local":
		logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}

func NewWriter(logPath string) (io.Writer, error) {
	if logPath == "" {
		return os.Stdout, nil
	}
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return io.MultiWriter(os.Stdout, logFile), nil
}
