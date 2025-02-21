-- name: GetStatusById :one
SELECT * FROM payment_status WHERE id = $1 LIMIT 1;

-- name: ListStatus :many
SELECT * FROM payment_status;

-- name: GetKindById :one
SELECT * FROM payment_kind WHERE id = $1 LIMIT 1;

-- name: GetPaymentKind :many
SELECT * FROM payment_kind;

-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1 LIMIT 1;

-- name: GetPaymentWithKindAndStatusName :one
SELECT 
    p.*,
    pk.name as kind_name,
    ps.name as status_name
FROM payments p
JOIN payment_kind pk ON p.kind_id = pk.id
JOIN payment_status ps ON p.status_id = ps.id
WHERE p.id = $1 
LIMIT 1;

-- name: GetPaymentByOrderId :one
SELECT * FROM payments WHERE orderID = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments ORDER BY created_at DESC;

-- name: CreatePayment :one
INSERT INTO payments (
    id,
    orderID,
    totalAmount,
    created_at,
    integratorExternalID,
    kind_id,
    status_id
    ) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdatePayment :exec
UPDATE payments SET orderID = $2, totalAmount = $3, created_at = $4, integratorExternalID = $5, kind_id = $6, status_id = $7 WHERE id = $1;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = $1;