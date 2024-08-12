package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type (
	ContextKey string
)

var (
	logKey         ContextKey = "logger"
	gcpSlogMapping            = map[string]string{
		"level": "severity",
		"msg":   "message",
	}
)

func isGCP() bool {
	return os.Getenv("K_SERVICE") != ""
}

func NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{}

	if isGCP() {
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if v, ok := gcpSlogMapping[a.Key]; ok {
				return slog.Attr{Key: v, Value: a.Value}
			}
			return a
		}
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := NewLogger()
		logger = logger.With("method", r.Method, "path", r.URL.Path)

		if isGCP() {
			logger = logger.With("logging.googleapis.com/trace", r.Header.Get("Traceparent"))
		}

		ctx := context.WithValue(r.Context(), logKey, logger)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(logKey).(*slog.Logger); ok {
		return logger
	}

	return NewLogger()
}
