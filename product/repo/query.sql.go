// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package repo

import (
	"context"

	"github.com/shopspring/decimal"
)

const getProduct = `-- name: GetProduct :one
SELECT 
    product.id, 
    product.category_id, 
    category.name AS category_name,
    product.name, 
    product.description, 
    product.price,
    product.discounted_price
FROM 
    product, 
    category
WHERE 
    product.id = ?1 AND
    category.id = product.category_id
GROUP BY category_id
`

type GetProductRow struct {
	ID              int64
	CategoryID      int64
	CategoryName    string
	Name            string
	Description     string
	Price           decimal.Decimal
	DiscountedPrice decimal.Decimal
}

func (q *Queries) GetProduct(ctx context.Context, id int64) (GetProductRow, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i GetProductRow
	err := row.Scan(
		&i.ID,
		&i.CategoryID,
		&i.CategoryName,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.DiscountedPrice,
	)
	return i, err
}

const getProductByCategory = `-- name: GetProductByCategory :many
SELECT 
    product.id, 
    product.category_id, 
    category.name AS category_name,
    product.name, 
    product.description, 
    product.price,
    product.discounted_price
FROM 
    product, 
    category
WHERE 
    product.category_id = ?1 AND
    category.id = product.category_id
GROUP BY category_id
`

type GetProductByCategoryRow struct {
	ID              int64
	CategoryID      int64
	CategoryName    string
	Name            string
	Description     string
	Price           decimal.Decimal
	DiscountedPrice decimal.Decimal
}

func (q *Queries) GetProductByCategory(ctx context.Context, categoryID int64) ([]GetProductByCategoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductByCategory, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductByCategoryRow
	for rows.Next() {
		var i GetProductByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.CategoryID,
			&i.CategoryName,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.DiscountedPrice,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProducts = `-- name: GetProducts :many
SELECT 
    product.id, 
    product.category_id, 
    category.name AS category_name,
    product.name, 
    product.description, 
    product.price,
    product.discounted_price
FROM
    product, 
    category
WHERE 
    category.id = product.category_id
GROUP BY category_id
`

type GetProductsRow struct {
	ID              int64
	CategoryID      int64
	CategoryName    string
	Name            string
	Description     string
	Price           decimal.Decimal
	DiscountedPrice decimal.Decimal
}

func (q *Queries) GetProducts(ctx context.Context) ([]GetProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsRow
	for rows.Next() {
		var i GetProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.CategoryID,
			&i.CategoryName,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.DiscountedPrice,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
