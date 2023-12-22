package apiserver

import (
	"log/slog"
	"time"
)

const (
	ContextCancelTimeout = 5 * time.Second
	ShutdownTimeout      = 2 * time.Second
	ServerReadTimeout    = 10 * time.Second
	ServerWriteTimeout   = 10 * time.Second
	ServerIdleTimeout    = 60 * time.Second

	WorkerCount = 100
	apiPrefix   = "/task"
)

type apiServer struct {
	logLevel slog.Level
	logger   *slog.Logger
}

type Option func(*apiServer)

func WithLogger(logger *slog.Logger) Option {
	return func(s *apiServer) {
		s.logger = logger
	}
}

func WithLogLevel(level string) Option {
	return func(s *apiServer) {
		var logLevel slog.Level

		switch level {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelInfo
		}
		s.logLevel = logLevel
	}
}
