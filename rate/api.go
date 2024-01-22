package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bojanz/currency"
	"github.com/shopspring/decimal"
)

type Client interface {
	Latest(ctx context.Context) (*Latest, error)
}

type Latest struct {
	Disclaimer string
	License    string
	Timestamp  time.Time
	Base       string
	Rates      map[string]decimal.Decimal
}

type UnsupportedSymbolError struct {
	Symbol string
}

type ClientError struct {
	StatusCode int
	Body       []byte
}

func (r *Latest) Convert(amount currency.Amount, symbol string) (currency.Amount, error) {
	if symbol == r.Base {
		return amount, nil
	}

	var zero currency.Amount

	rate, ok := r.Rates[symbol]
	if !ok {
		return zero, UnsupportedSymbolError{Symbol: symbol}
	}

	n, err := decimal.NewFromString(amount.Number())
	if err != nil {
		return zero, err
	}

	n = n.Mul(rate)
	converted, err := currency.NewAmount(n.String(), symbol)
	if err != nil {
		return zero, err
	}

	return converted.Round(), nil
}

func (r *Latest) UnmarshalJSON(data []byte) error {
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

	*r = Latest{
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

func (e ClientError) Error() string {
	return fmt.Sprintf("rate: response status %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}
