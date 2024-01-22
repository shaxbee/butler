-- name: GetProduct :one
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
    product.id = sqlc.arg(id) AND
    category.id = product.category_id
GROUP BY category_id;

-- name: GetProducts :many
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
GROUP BY category_id;

-- name: GetProductByCategory :many
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
    product.category_id = sqlc.arg(category_id) AND
    category.id = product.category_id
GROUP BY category_id;