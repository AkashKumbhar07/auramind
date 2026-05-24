package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		window:   window,
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		clientIP := extractClientIP(ctx)

		rl.mu.Lock()
		count := rl.requests[clientIP]
		if count >= rl.limit {
			rl.mu.Unlock()
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		rl.requests[clientIP] = count + 1
		rl.mu.Unlock()

		return handler(ctx, req)
	}
}

func extractClientIP(ctx context.Context) string {
	// Extract from metadata/context
	return "unknown"
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	for range ticker.C {
		rl.mu.Lock()
		rl.requests = make(map[string]int)
		rl.mu.Unlock()
	}
}
