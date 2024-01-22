package rate

import (
	"net/url"
	"testing"

	"github.com/bojanz/currency"
	"github.com/shopspring/decimal"
)

// Client is the interface implemented by types that can retrieve rates.
var (
	_ Client = (*RestClient)(nil)
	_ Client = (*SyntheticClient)(nil)
)

func verifyLatest(t testing.TB, base string, symbols []string, latest *Latest) {
	t.Helper()

	if latest == nil {
		t.Fatal("expected latest rates")
	}

	if expected := "https://openexchangerates.org/license"; latest.License != expected {
		t.Errorf("expected license %q, got %q", expected, latest.License)
	}

	if latest.Disclaimer == "" {
		t.Error("expected disclaimer")
	}

	if latest.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	if latest.Base != base {
		t.Errorf("expected base %s, got %s", base, latest.Base)
	}

	verifyRates(t, symbols, latest.Rates)
}

func verifyRates(t testing.TB, symbols []string, rates map[string]decimal.Decimal) {
	t.Helper()

	for _, symbol := range symbols {
		rate, ok := rates[symbol]
		switch {
		case !ok:
			t.Errorf("expected symbol %q in rates", symbol)
		case rate.IsZero():
			t.Errorf("expected non-zero rate for %q", symbol)
		}
	}
}

func mustConvert(t testing.TB, latest *Latest, amount currency.Amount, symbol string) currency.Amount {
	t.Helper()

	converted, err := latest.Convert(amount, symbol)
	if err != nil {
		t.Fatal("convert:", err)
	}

	return converted
}

func dec(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func amount(n string, symbol string) currency.Amount {
	a, err := currency.NewAmount(n, symbol)
	if err != nil {
		panic(err)
	}
	return a
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
