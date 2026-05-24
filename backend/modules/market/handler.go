package market

import (
	"context"

	"github.com/AkashKumbhar07/auramind/backend/framework/validation"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	svc    *Service
	logger *zap.Logger
}

func NewHandler(svc *Service, logger *zap.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) GetPairs(ctx context.Context, req *GetPairsRequest) (*GetPairsResponse, error) {
	if err := validation.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	resp, err := h.svc.GetPairs(ctx, req.Exchange)
	if err != nil {
		h.logger.Error("get pairs failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "get pairs: %v", err)
	}

	return resp, nil
}

func (h *Handler) GetCandles(ctx context.Context, req *GetCandlesRequest) (*GetCandlesResponse, error) {
	if err := validation.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	resp, err := h.svc.GetCandles(ctx, req)
	if err != nil {
		h.logger.Error("get candles failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "get candles: %v", err)
	}

	return resp, nil
}

func (h *Handler) Analyze(ctx context.Context, req *AnalyzeRequest) (*AnalyzeResponse, error) {
	if err := validation.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	resp, err := h.svc.Analyze(ctx, req)
	if err != nil {
		h.logger.Error("analyze failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "analyze: %v", err)
	}

	return resp, nil
}
