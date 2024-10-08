package logger

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, message string)
	Error(ctx context.Context, err error)
	Info(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
}
