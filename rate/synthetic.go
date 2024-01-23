package rate

import (
	"context"
	"slices"

	"github.com/shopspring/decimal"
)

type SyntheticClient struct {
	delegate Client
	base     string
}

func NewSyntheticClient(delegate Client, base string) *SyntheticClient {
	return &SyntheticClient{
		delegate: delegate,
		base:     base,
	}
}

func (c *SyntheticClient) Latest(ctx context.Context, req LatestRequest) (*LatestResponse, error) {
	req.Base = DefaultBase

	if !slices.Contains(req.Symbols, c.base) {
		req.Symbols = append(req.Symbols, c.base)
	}

	latest, err := c.delegate.Latest(ctx, req)
	if err != nil {
		return nil, err
	}

	rate := decimal.NewFromInt(1).Div(latest.Rates[c.base])
	synthetic := map[string]decimal.Decimal{
		DefaultBase: rate,
	}
	for symbol, rate := range latest.Rates {
		if symbol == req.Base {
			continue
		}

		synthetic[symbol] = rate.Mul(rate)
	}

	return &LatestResponse{
		ETag:       latest.ETag,
		Disclaimer: latest.Disclaimer,
		License:    latest.License,
		Timestamp:  latest.Timestamp,
		Base:       c.base,
		Rates:      synthetic,
	}, nil
}
