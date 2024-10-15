package cache

import "context"

type Cache interface {
	HealthCheck(ctx context.Context) error
	AddData(ctx context.Context, args ...any) (string, error)
	GetData(ctx context.Context, args ...any) (any, error)
	GetKeys(ctx context.Context, pattern string) ([]string, error)
	DeleteData(ctx context.Context, args ...any) (any, error)
}
