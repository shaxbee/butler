package rate

import (
	"context"

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

func (c *SyntheticClient) Latest(ctx context.Context) (*Latest, error) {
	latest, err := c.delegate.Latest(ctx)
	if err != nil {
		return nil, err
	}

	base := decimal.NewFromInt(1).Div(latest.Rates[c.base])
	synthetic := map[string]decimal.Decimal{
		DefaultBase: base,
	}
	for symbol, rate := range latest.Rates {
		if symbol == c.base {
			continue
		}

		synthetic[symbol] = rate.Mul(base)
	}

	return &Latest{
		Disclaimer: latest.Disclaimer,
		License:    latest.License,
		Timestamp:  latest.Timestamp,
		Base:       c.base,
		Rates:      synthetic,
	}, nil
}
