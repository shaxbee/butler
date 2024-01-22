package rate

import (
	"github.com/bojanz/currency"
	"github.com/shopspring/decimal"
)

type Converter interface {
	Convert(amount decimal.Decimal, pair CurrencyPair) (currency.Amount, error)
}

type CurrencyPair struct {
	From string
	To   string
}
