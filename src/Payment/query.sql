-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1 LIMIT 1;

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
    integratorExternalID
    ) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdatePayment :exec
UPDATE payments SET orderID = $2, totalAmount = $3, created_at = $4, integratorExternalID = $5 WHERE id = $1;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = $1;