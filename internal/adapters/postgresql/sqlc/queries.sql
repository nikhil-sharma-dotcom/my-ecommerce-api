-- name: CreateProduct :one
INSERT INTO products (name, price_in_cents, quantity)
VALUES ($1, $2, $3)
RETURNING id, name, price_in_cents, quantity, created_at;

-- name: GetProduct :one
SELECT id, name, price_in_cents, quantity, created_at
FROM products
WHERE id = $1;

-- name: GetProductForUpdate :one
SELECT id, name, price_in_cents, quantity, created_at
FROM products
WHERE id = $1
FOR UPDATE;

-- name: ListProducts :many
SELECT id, name, price_in_cents, quantity, created_at
FROM products
ORDER BY created_at DESC;

-- name: UpdateProduct :exec
UPDATE products
SET name = $1, price_in_cents = $2, quantity = $3
WHERE id = $4;

-- name: UpdateProductQuantity :exec
UPDATE products
SET quantity = quantity - $1
WHERE id = $2 AND quantity >= $1;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (customer_id)
VALUES ($1)
RETURNING id, customer_id, created_at;

-- name: GetOrder :one
SELECT id, customer_id, created_at
FROM orders
WHERE id = $1;

-- name: ListOrdersByCustomer :many
SELECT id, customer_id, created_at
FROM orders
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents)
VALUES ($1, $2, $3, $4)
RETURNING id, order_id, product_id, quantity, price_cents;

-- name: GetOrderItems :many
SELECT id, order_id, product_id, quantity, price_cents
FROM order_items
WHERE order_id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, password_hash, first_name, last_name, role, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, first_name, last_name, role, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password_hash, first_name, last_name, role, created_at
FROM users
WHERE id = $1;
