package cart

import (
	"time"

	"github.com/shopspring/decimal"
)

type Cart struct {
	Items []Item
}

func (c Cart) Total() decimal.Decimal {
	total := decimal.Zero
	for _, item := range c.Items {
		total = total.Add(item.Total)
	}

	return total
}

type Item struct {
	Created         time.Time
	ProductID       string
	Name            string
	Price           decimal.Decimal
	DiscountedPrice decimal.Decimal
	Quantity        decimal.Decimal
	Total           decimal.Decimal
}
