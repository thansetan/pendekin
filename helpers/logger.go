package helpers

import (
	"log/slog"
	"os"
)

func NewLogger(mode string) *slog.Logger {
	var logHandler slog.Handler
	opts := &slog.HandlerOptions{
		AddSource: true,
	}

	switch mode {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stderr, opts)
	case "text":
		logHandler = slog.NewTextHandler(os.Stderr, opts)
	}
	logger := slog.New(logHandler)

	return logger
}
