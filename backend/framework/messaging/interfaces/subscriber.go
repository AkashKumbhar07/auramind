package interfaces

import "context"

type Handler func(ctx context.Context, subject string, data []byte)

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler Handler) error
	Close() error
}
