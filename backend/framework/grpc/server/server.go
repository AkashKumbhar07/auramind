package server

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server
	port    int
	logger  *zap.Logger
	health  *health.Server
}

func New(port int, logger *zap.Logger, opts ...grpc.ServerOption) *Server {
	gs := grpc.NewServer(opts...)

	s := &Server{
		Server:  gs,
		port:    port,
		logger:  logger,
		health:  health.NewServer(),
	}

	grpc_health_v1.RegisterHealthServer(gs, s.health)
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	reflection.Register(gs)

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	s.logger.Info("gRPC server starting", zap.String("addr", addr))
	return s.Serve(lis)
}

func (s *Server) Shutdown() {
	s.logger.Info("gRPC server shutting down")
	s.GracefulStop()
}
