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
		logHandler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(logHandler)

	return logger
}
