package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

func New(levelStr string) *slog.Logger {
	return newLogger(levelStr, os.Stdout)
}

func newLogger(levelStr string, w io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: parseLevel(levelStr),
	})
	return slog.New(handler)
}

func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
