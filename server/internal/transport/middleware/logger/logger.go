package logger

import (
	"context"
	"log/slog"
	"net/http"
)

const key = "logger"

func Middleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(key).(*slog.Logger)
	if !ok {
		return nil
	}
	return logger
}
