package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const fieldsKey contextKey = "logger_fields"

func Init(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	})

	slog.SetDefault(slog.New(handler))
}

func WithFields(ctx context.Context, fields ...any) context.Context {
	existing := fieldsFromContext(ctx)
	merged := make([]any, 0, len(existing)+len(fields))
	merged = append(merged, existing...)
	merged = append(merged, fields...)
	return context.WithValue(ctx, fieldsKey, merged)
}

func Info(ctx context.Context, msg string, args ...any) {
	allArgs := append(fieldsFromContext(ctx), args...)
	slog.InfoContext(ctx, msg, allArgs...)
}

func Error(ctx context.Context, msg string, args ...any) {
	allArgs := append(fieldsFromContext(ctx), args...)
	slog.ErrorContext(ctx, msg, allArgs...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	allArgs := append(fieldsFromContext(ctx), args...)
	slog.WarnContext(ctx, msg, allArgs...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	allArgs := append(fieldsFromContext(ctx), args...)
	slog.DebugContext(ctx, msg, allArgs...)
}

func fieldsFromContext(ctx context.Context) []any {
	if fields, ok := ctx.Value(fieldsKey).([]any); ok {
		return fields
	}
	return nil
}
