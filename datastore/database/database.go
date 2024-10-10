package database

import "context"

type Database interface {
	CloseDB(ctx context.Context, cancel context.CancelFunc)
	HealthCheck(ctx context.Context) error
	AddData(ctx context.Context, args ...any) (string, error)
	AddMultipleData(ctx context.Context, args ...any) ([]string, error)
	GetData(ctx context.Context, args ...any) (any, error)
	GetMultipleData(ctx context.Context, args ...any) ([]any, error)
	UpdateData(ctx context.Context, args ...any) (any, error)
	UpdateMultipleData(ctx context.Context, args ...any) (any, error)
	DeleteData(ctx context.Context, args ...any) (any, error)
	DeleteMultipleData(ctx context.Context, args ...any) (any, error)
}
