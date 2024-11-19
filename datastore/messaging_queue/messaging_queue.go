package messagingqueue

import "context"

type MessageQueue interface {
	HealthCheck(ctx context.Context) error
	ProduceMessage(ctx context.Context, args ...any) error
	ConsumeMessage(ctx context.Context, args ...any) (any, error)
}
