package product

import (
	"github.com/bojanz/currency"
)

type CategoryProducts struct {
	Category string
	Products []Product
}

func (c CategoryProducts) IsZero() bool {
	return c.Category == "" && len(c.Products) == 0
}

type Product struct {
	ID              string
	Category        string
	Name            string
	Price           currency.Amount
	DiscountedPrice currency.Amount
	Image           string
	Description     string
}

func GroupByCategory(products []Product) []CategoryProducts {
	var res []CategoryProducts
	for _, product := range products {
		var found bool
		for i, c := range res {
			if c.Category == product.Category {
				res[i].Products = append(c.Products, product)
				found = true
				break
			}
		}

		if !found {
			res = append(res, CategoryProducts{
				Category: product.Category,
				Products: []Product{product},
			})
		}
	}

	return res
}
