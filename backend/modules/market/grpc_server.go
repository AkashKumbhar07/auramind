package market

import (
	"context"
	"log"

	pb "github.com/AkashKumbhar07/auramind/backend/framework/grpc/generated/market"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedMarketServiceServer
	svc *Service
}

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) GetPairs(ctx context.Context, req *pb.GetPairsRequest) (*pb.GetPairsResponse, error) {
	domainResp, err := s.svc.GetPairs(ctx, req.GetExchange())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get pairs: %v", err)
	}

	pairs := make([]*pb.MarketPair, len(domainResp.Pairs))
	for i, p := range domainResp.Pairs {
		pairs[i] = &pb.MarketPair{
			Id:         p.ID,
			Symbol:     p.Symbol,
			BaseAsset:  p.BaseAsset,
			QuoteAsset: p.QuoteAsset,
			Exchange:   p.Exchange,
			Active:     p.Active,
		}
	}

	return &pb.GetPairsResponse{Pairs: pairs}, nil
}

func (s *GRPCServer) GetCandles(ctx context.Context, req *pb.GetCandlesRequest) (*pb.GetCandlesResponse, error) {
	domainReq := &GetCandlesRequest{
		Symbol:   req.GetSymbol(),
		Exchange: req.GetExchange(),
		Interval: req.GetInterval(),
		Limit:    int(req.GetLimit()),
	}

	domainResp, err := s.svc.GetCandles(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get candles: %v", err)
	}

	candles := make([]*pb.Candle, len(domainResp.Candles))
	for i, c := range domainResp.Candles {
		candles[i] = &pb.Candle{
			Symbol:    c.Symbol,
			Exchange:  c.Exchange,
			Interval:  c.Interval,
			Open:      c.Open,
			High:      c.High,
			Low:       c.Low,
			Close:     c.Close,
			Volume:    c.Volume,
			Timestamp: c.Timestamp,
		}
	}

	return &pb.GetCandlesResponse{Candles: candles}, nil
}

func (s *GRPCServer) GetTicker(ctx context.Context, req *pb.GetTickerRequest) (*pb.GetTickerResponse, error) {
	_ = req
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *GRPCServer) SubscribePrices(req *pb.SubscribePricesRequest, stream pb.MarketService_SubscribePricesServer) error {
	_ = req
	_ = stream
	log.Println("SubscribePrices: stream endpoint not yet wired")
	return nil
}
