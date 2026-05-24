package telemetry

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/processors/minsev"
)

func NewLogger(env string, logLevel minsev.Severity) *slog.Logger {
	level := newLogLevelFromMinSev(logLevel)

	stdoutHandler := tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  true,
		TimeFormat: time.Kitchen,
		Level:      level,
	})

	if env == "production" {
		stdoutHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
	}

	otelHandler := otelslog.NewHandler("joke-service", otelslog.WithSource(true))

	handler := slog.NewMultiHandler(stdoutHandler, otelHandler)

	return slog.New(handler)
}

func newLogLevelFromMinSev(severity minsev.Severity) slog.Level {
	switch severity {
	case minsev.SeverityError:
		return slog.LevelError
	case minsev.SeverityWarn:
		return slog.LevelWarn
	case minsev.SeverityInfo:
		return slog.LevelInfo
	case minsev.SeverityDebug:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
