package market

import (
	"context"
	"fmt"
	"math"

	coreIndicators "github.com/AkashKumbhar07/auramind/backend/core/indicators"
	"github.com/AkashKumbhar07/auramind/backend/framework/errors"
	"go.uber.org/zap"
)

type Service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) GetPairs(ctx context.Context, exchange string) (*GetPairsResponse, error) {
	pairs, err := s.repo.GetPairs(ctx, exchange)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "fetch pairs", err)
	}
	return &GetPairsResponse{Pairs: pairs}, nil
}

func (s *Service) GetCandles(ctx context.Context, req *GetCandlesRequest) (*GetCandlesResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 100
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	candles, err := s.repo.GetCandles(ctx, req.Symbol, req.Exchange, req.Interval, req.Limit)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "fetch candles", err)
	}

	return &GetCandlesResponse{Candles: candles}, nil
}

func (s *Service) Analyze(ctx context.Context, req *AnalyzeRequest) (*AnalyzeResponse, error) {
	if len(req.Candles) < 30 {
		return nil, errors.BadRequest("need at least 30 candles for analysis")
	}

	prices := make([]float64, len(req.Candles))
	for i, c := range req.Candles {
		prices[i] = c.Close
	}

	indicators := s.calcIndicators(prices)

	score := s.calcBullishScore(prices, indicators)
	riskLevel := s.calcRiskLevel(score, indicators)
	stopLoss := s.calcStopLoss(req.Candles)
	reasons := s.generateReasons(score, indicators, riskLevel)
	invalidation := s.calcInvalidation(prices, indicators, riskLevel)

	return &AnalyzeResponse{
		Symbol:             req.Symbol,
		BullishProbability: math.Round(score*100) / 100,
		ConfidenceScore:    math.Round(s.calcConfidence(indicators)*100) / 100,
		RiskLevel:          riskLevel,
		SuggestedStopLoss:  math.Round(stopLoss*100) / 100,
		Indicators:         indicators,
		Reasons:            reasons,
		Invalidation:       invalidation,
	}, nil
}

func (s *Service) calcIndicators(prices []float64) IndicatorResult {
	var result IndicatorResult

	// RSI (14)
	rsiVals := coreIndicators.RSIValues(prices, 14)
	if len(rsiVals) > 0 {
		v := math.Round(rsiVals[len(rsiVals)-1]*100) / 100
		result.RSI = &v
	}

	// MACD (12, 26, 9)
	macdVals := coreIndicators.MACDValues(prices, 12, 26, 9)
	if len(macdVals) > 0 {
		last := macdVals[len(macdVals)-1]
		result.MACD = &MACDValue{
			Value:     math.Round(last.Value*100) / 100,
			Signal:    math.Round(last.Signal*100) / 100,
			Histogram: math.Round(last.Histogram*100) / 100,
		}
	}

	// Bollinger Bands (20, 2)
	bbVals := coreIndicators.BollingerValues(prices, 20, 2)
	if len(bbVals) > 0 {
		last := bbVals[len(bbVals)-1]
		result.BB = &BBValue{
			Upper:  math.Round(last.Upper*100) / 100,
			Middle: math.Round(last.Middle*100) / 100,
			Lower:  math.Round(last.Lower*100) / 100,
		}
	}

	// EMA (9)
	ema9 := coreIndicators.EMAValues(prices, 9)
	if len(ema9) > 0 {
		v := math.Round(ema9[len(ema9)-1]*100) / 100
		result.EMA9 = &v
	}

	// EMA (21)
	ema21 := coreIndicators.EMAValues(prices, 21)
	if len(ema21) > 0 {
		v := math.Round(ema21[len(ema21)-1]*100) / 100
		result.EMA21 = &v
	}

	return result
}

func (s *Service) calcBullishScore(prices []float64, ind IndicatorResult) float64 {
	score := 50.0
	lastPrice := prices[len(prices)-1]

	if ind.RSI != nil {
		if *ind.RSI < 30 {
			score += 15 // oversold bounce potential
		} else if *ind.RSI > 70 {
			score -= 15 // overbought
		} else if *ind.RSI > 50 {
			score += 5
		} else {
			score -= 5
		}
	}

	if ind.MACD != nil {
		if ind.MACD.Histogram > 0 {
			score += 10
		} else {
			score -= 10
		}
		if ind.MACD.Value > ind.MACD.Signal {
			score += 5
		} else {
			score -= 5
		}
	}

	if ind.BB != nil {
		if lastPrice <= ind.BB.Lower*1.02 {
			score += 10 // near lower band = support
		} else if lastPrice >= ind.BB.Upper*0.98 {
			score -= 10 // near upper band = resistance
		}
	}

	if ind.EMA9 != nil && ind.EMA21 != nil {
		if *ind.EMA9 > *ind.EMA21 {
			score += 10 // bullish cross
		} else {
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

func (s *Service) calcConfidence(ind IndicatorResult) float64 {
	count := 0
	if ind.RSI != nil {
		count++
	}
	if ind.MACD != nil {
		count++
	}
	if ind.BB != nil {
		count++
	}
	if ind.EMA9 != nil && ind.EMA21 != nil {
		count++
	}

	if count == 0 {
		return 0
	}
	return float64(count) / 4.0 * 100
}

func (s *Service) calcRiskLevel(score float64, ind IndicatorResult) string {
	volatility := 0.0
	if ind.BB != nil && ind.BB.Middle > 0 {
		volatility = (ind.BB.Upper - ind.BB.Lower) / ind.BB.Middle
	}

	if volatility > 0.1 {
		return "high"
	}
	if score < 30 || score > 70 {
		return "medium"
	}
	return "low"
}

func (s *Service) calcStopLoss(candles []Candle) float64 {
	if len(candles) < 10 {
		return 3.0
	}

	min := candles[0].Close
	max := candles[0].Close
	for _, c := range candles {
		if c.Close < min {
			min = c.Close
		}
		if c.Close > max {
			max = c.Close
		}
	}

	range_ := max - min
	avg := (max + min) / 2
	if avg == 0 {
		return 3.0
	}

	pct := (range_ / avg) * 100
	if pct < 1 {
		return 1.5
	}
	if pct > 10 {
		return 5.0
	}
	return math.Round(pct*10) / 10
}

func (s *Service) generateReasons(score float64, ind IndicatorResult, risk string) []string {
	var reasons []string

	if ind.RSI != nil {
		if *ind.RSI < 30 {
			reasons = append(reasons, fmt.Sprintf("RSI at %.1f — oversold territory, potential bounce", *ind.RSI))
		} else if *ind.RSI > 70 {
			reasons = append(reasons, fmt.Sprintf("RSI at %.1f — overbought, caution advised", *ind.RSI))
		} else if *ind.RSI > 50 {
			reasons = append(reasons, fmt.Sprintf("RSI at %.1f — bullish momentum", *ind.RSI))
		} else {
			reasons = append(reasons, fmt.Sprintf("RSI at %.1f — bearish pressure", *ind.RSI))
		}
	}

	if ind.MACD != nil {
		if ind.MACD.Histogram > 0 {
			reasons = append(reasons, fmt.Sprintf("MACD histogram positive at %.2f — bullish crossover", ind.MACD.Histogram))
		} else {
			reasons = append(reasons, fmt.Sprintf("MACD histogram negative at %.2f — bearish crossover", ind.MACD.Histogram))
		}
	}

	if ind.BB != nil {
		reasons = append(reasons, fmt.Sprintf("Bollinger Bands: upper=%.2f, lower=%.2f", ind.BB.Upper, ind.BB.Lower))
	}

	if ind.EMA9 != nil && ind.EMA21 != nil {
		if *ind.EMA9 > *ind.EMA21 {
			reasons = append(reasons, fmt.Sprintf("EMA9 (%.2f) above EMA21 (%.2f) — bullish alignment", *ind.EMA9, *ind.EMA21))
		} else {
			reasons = append(reasons, fmt.Sprintf("EMA9 (%.2f) below EMA21 (%.2f) — bearish alignment", *ind.EMA9, *ind.EMA21))
		}
	}

	reasons = append(reasons, fmt.Sprintf("Risk level: %s", risk))

	return reasons
}

func (s *Service) calcInvalidation(prices []float64, ind IndicatorResult, risk string) string {
	if len(prices) < 5 {
		return "insufficient data"
	}

	support := prices[0]
	for _, p := range prices {
		if p < support {
			support = p
		}
	}

	if ind.BB != nil && ind.BB.Lower < support {
		support = ind.BB.Lower
	}

	return fmt.Sprintf("Break below support level at %.2f", support)
}

func (s *Service) SubmitCandle(ctx context.Context, candle *Candle) error {
	if err := s.repo.SaveCandle(ctx, candle); err != nil {
		return errors.Wrap(errors.KindInternal, "save candle", err)
	}
	return nil
}
