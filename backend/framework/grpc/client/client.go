package client

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	logger *zap.Logger
}

func New(host string, port int, logger *zap.Logger, opts ...grpc.DialOption) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(10 * time.Second),
	}

	allOpts := append(defaultOpts, opts...)

	conn, err := grpc.NewClient(addr, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", addr, err)
	}

	logger.Info("gRPC client connected", zap.String("addr", addr))

	return &Client{
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Close() error {
	c.logger.Info("gRPC client closing")
	return c.conn.Close()
}
