package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(),
	}))
	slog.SetDefault(logger)
	return logger
}

func getLogLevel() slog.Level {
	level := env.GetString("LOG_LEVEL", "INFO")

	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type AsynqLoggerAdapter struct {
	logger *slog.Logger
}

func NewAsynqLoggerAdapter(l *slog.Logger) *AsynqLoggerAdapter {
	return &AsynqLoggerAdapter{logger: l}
}

func (a *AsynqLoggerAdapter) Debug(args ...interface{}) {
	a.logger.Debug("asynq debug", "details", args)
}

func (a *AsynqLoggerAdapter) Info(args ...interface{}) {
	a.logger.Info("asynq info", "details", args)
}

func (a *AsynqLoggerAdapter) Warn(args ...interface{}) {
	a.logger.Warn("asynq warn", "details", args)
}

func (a *AsynqLoggerAdapter) Error(args ...interface{}) {
	a.logger.Error("asynq error", "details", args)
}

func (a *AsynqLoggerAdapter) Fatal(args ...interface{}) {
	a.logger.Error("asynq fatal error", "details", args)
	os.Exit(1)
}
