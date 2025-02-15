-- name: GetStatus :one
SELECT * FROM status WHERE id = $1 LIMIT 1;

-- name: ListStatuses :many
SELECT * FROM status ORDER BY status_name;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders ORDER BY created_at DESC;

-- name: CreateOrder :one
INSERT INTO orders (id, customer_id, status_id, created_at, total_price, discount) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateOrder :exec
UPDATE orders SET status_id = $2, total_price = $3, discount = $4 WHERE id = $1;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: GetOrderWithItems :many
SELECT o.id, o.customer_id, o.status_id, o.created_at, o.total_price, o.discount, i.id AS item_id, i.price
FROM orders o
JOIN order_items oi ON o.id = oi.order_id
JOIN items i ON i.id = oi.item_id
WHERE o.id = $1;