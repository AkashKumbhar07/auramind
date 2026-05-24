package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logging(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		dur := time.Since(start)
		code := status.Code(err)

		logger.Info("gRPC call",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", dur),
			zap.String("code", code.String()),
		)

		return resp, err
	}
}

func LoggingStream(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, stream)

		dur := time.Since(start)
		logger.Info("gRPC stream",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", dur),
			zap.Bool("is_server_stream", info.IsServerStream),
			zap.Bool("is_client_stream", info.IsClientStream),
		)

		return err
	}
}
