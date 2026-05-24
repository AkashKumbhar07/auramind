package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	msginterfaces "github.com/AkashKumbhar07/auramind/backend/framework/messaging/interfaces"
)

type Subscriber struct {
	conn *nats.Conn
}

func NewSubscriber(url string) (*Subscriber, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}
	return &Subscriber{conn: conn}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, subject string, handler msginterfaces.Handler) error {
	_, err := s.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(ctx, msg.Subject, msg.Data)
	})
	if err != nil {
		return fmt.Errorf("subscribe to %s: %w", subject, err)
	}
	return nil
}

func (s *Subscriber) Close() error {
	s.conn.Close()
	return nil
}

// Ensure interfaces are satisfied
var _ msginterfaces.Publisher = (*Publisher)(nil)
var _ msginterfaces.Subscriber = (*Subscriber)(nil)
