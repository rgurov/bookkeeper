package logging

import (
	"context"

	"golang.org/x/exp/slog"
)

type ctxLogger struct{}

func LoggerWithContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}
