-- name: CreateUser :one
INSERT INTO users (
  id, email, full_name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByID :one
SELECT sqlc.embed(users)
 FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
  set email = $2,
  full_name = $3
WHERE id = $1
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (
  name, description, category, base_price, is_customizable, customization_fields
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetProducts :many
SELECT * FROM products
ORDER BY created_at DESC;

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = $1;

-- name: CreateCustomer :one
INSERT INTO customers (
  email, phone, name, company_name, address, city, postal_code, country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateCustomer :exec
UPDATE customers
SET email = $2, phone = $3, name = $4, company_name = $5, address = $6, city = $7, postal_code = $8, country = $9, updated_at = NOW()
WHERE id = $1;

-- name: GetAllCustomers :many
SELECT * FROM customers ORDER BY created_at DESC;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1 ORDER BY created_at DESC;

-- name: CreateOrder :one
INSERT INTO orders (
  order_number, customer_id, order_type, status, total_price, notes, order_date, expected_completion_date, actual_completion_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateOrder :exec
UPDATE orders
SET customer_id = $2, order_type = $3, status= $4, total_price = $5, notes = $6, order_date = $7, expected_completion_date = $8, actual_completion_date = $9
WHERE id = $1;

-- name: GetOrders :many
SELECT *, count(*) OVER() as total_count
FROM orders;

-- name: CreateOrderItems :one
INSERT INTO order_items (
  order_id, product_id, quantity, price_per_unit
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items
WHERE order_id = $1;

-- name: GetOrderWithItems :one
SELECT
  o.id,
  o.order_number,
  o.customer_id,
  o.order_type,
  o.status,
  o.total_price,
  o.notes,
  o.order_date,
  o.expected_completion_date,
  o.actual_completion_date,
  COALESCE(
    json_agg(
      json_build_object(
        'item_id', oi.id,
        'product_id', oi.product_id,
        'quantity', oi.quantity,
        'price_per_unit', oi.price_per_unit
      )
    ) FILTER (WHERE oi.id IS NOT NULL),
    '[]'::json
  )::json AS items
FROM orders o
LEFT JOIN order_items oi ON o.id = oi.order_id
WHERE o.id = $1 GROUP BY o.id;

-- name: GetMaterials :many
SELECT * FROM materials ORDER BY created_at DESC;

-- name: CreateMaterial :one
INSERT INTO materials (
  name, thickness_mm, grade, description, cost_per_kg, stock_qty_kg, reorder_level_kg
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, name, thickness_mm, grade, description, cost_per_kg, stock_qty_kg, reorder_level_kg, created_at, updated_at;

-- name: GetMaterialByID :one
SELECT id, name, thickness_mm, grade, description, cost_per_kg, stock_qty_kg, reorder_level_kg, created_at, updated_at
FROM materials
WHERE id = $1 LIMIT 1;

-- name: ListMaterials :many
SELECT id, name, thickness_mm, grade, description, cost_per_kg, stock_qty_kg, reorder_level_kg, created_at, updated_at
FROM materials
ORDER BY created_at DESC;

-- name: UpdateMaterialStock :one
UPDATE materials
SET
  stock_qty_kg = stock_qty_kg + $2,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, stock_qty_kg, updated_at;