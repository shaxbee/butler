package rate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type Client interface {
	Latest(ctx context.Context, req LatestRequest) (*LatestResponse, error)
}

type LatestRequest struct {
	ETag            string
	Modified        time.Time
	Base            string
	Symbols         []string
	ShowAlternative bool
}

type LatestResponse struct {
	ETag       string
	Disclaimer string
	License    string
	Timestamp  time.Time
	Base       string
	Rates      map[string]decimal.Decimal
}

func (r *LatestResponse) Equal(o *LatestResponse) bool {
	if r == nil || o == nil {
		return r == o
	}

	if r.ETag != o.ETag {
		return false
	}

	if r.Disclaimer != o.Disclaimer {
		return false
	}

	if r.License != o.License {
		return false
	}

	if !r.Timestamp.Equal(o.Timestamp) {
		return false
	}

	if r.Base != o.Base {
		return false
	}

	return maps.Equal(r.Rates, o.Rates)
}

type UnsupportedSymbolError struct {
	Symbol string
}

type ClientError struct {
	StatusCode int
	Body       []byte
}

func (r *LatestResponse) Convert(amount decimal.Decimal, symbol string) (decimal.Decimal, error) {
	if symbol == r.Base {
		return amount, nil
	}

	var zero decimal.Decimal

	rate, ok := r.Rates[symbol]
	if !ok {
		return zero, UnsupportedSymbolError{Symbol: symbol}
	}

	return amount.Mul(rate), nil
}

func (r *LatestResponse) UnmarshalJSON(data []byte) error {
	src := struct {
		Disclaimer string                     `json:"disclaimer"`
		License    string                     `json:"license"`
		Timestamp  int64                      `json:"timestamp"`
		Base       string                     `json:"base"`
		Rates      map[string]decimal.Decimal `json:"rates"`
	}{}
	if err := json.Unmarshal(data, &src); err != nil {
		return err
	}

	*r = LatestResponse{
		ETag:       r.ETag,
		Disclaimer: src.Disclaimer,
		License:    src.License,
		Timestamp:  time.Unix(src.Timestamp, 0),
		Base:       src.Base,
		Rates:      src.Rates,
	}

	return nil
}

func (e UnsupportedSymbolError) Error() string {
	return "rate: unsupported symbol " + e.Symbol
}

func StatusCode(err error) int {
	var ce ClientError
	if !errors.As(err, &ce) {
		return 0
	}

	return ce.StatusCode
}

func (e ClientError) Error() string {
	return fmt.Sprintf("response status %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}
