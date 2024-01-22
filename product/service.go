package product

import (
	"context"

	"github.com/shaxbee/butler/product/repo"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID              int64
	CategoryID      int64
	CategoryName    string
	Name            string
	Description     string
	Price           decimal.Decimal
	DiscountedPrice decimal.Decimal
}

type Service struct {
	db repo.DBTX
}

func NewService(db repo.DBTX) *Service {
	return &Service{db: db}
}

func (s *Service) List(ctx context.Context, currencyCode string) ([]Product, error) {
	q := repo.New(s.db)

	products, err := q.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Product, len(products))
	for i, p := range products {
		result[i] = Product(p)
	}

	return result, nil
}
