-- name: CreatePayment :one
INSERT INTO payments (id, reservation_id, gateway, amount, currency, status, transaction_ref, payload)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1;

-- name: GetPaymentsByReservationID :many
SELECT * FROM payments WHERE reservation_id = $1 ORDER BY created_at DESC;

-- name: GetPaymentByTransactionRef :one
SELECT * FROM payments WHERE transaction_ref = $1;

-- name: ListPayments :many
SELECT p.*, r.pickup, r.destination 
FROM payments p
LEFT JOIN reservations r ON p.reservation_id = r.id
WHERE ($1::text IS NULL OR p.reservation_id ILIKE '%' || $1 || '%' OR p.transaction_ref ILIKE '%' || $1 || '%')
  AND ($2::payment_status IS NULL OR p.status = $2)
  AND ($3::payment_gateway IS NULL OR p.gateway = $3)
ORDER BY 
  CASE WHEN $4 = 'created_at' THEN p.created_at END DESC,
  CASE WHEN $4 = 'amount' THEN p.amount END DESC
LIMIT $5 OFFSET $6;

-- name: CountPayments :one
SELECT COUNT(*) FROM payments p
WHERE ($1::text IS NULL OR p.reservation_id ILIKE '%' || $1 || '%' OR p.transaction_ref ILIKE '%' || $1 || '%')
  AND ($2::payment_status IS NULL OR p.status = $2)
  AND ($3::payment_gateway IS NULL OR p.gateway = $3);

-- name: UpdatePayment :one
UPDATE payments
SET status = COALESCE($2, status),
    transaction_ref = COALESCE($3, transaction_ref),
    payload = COALESCE($4, payload)
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $2
WHERE id = $1
RETURNING *;

-- name: GetPaymentsByStatus :many
SELECT * FROM payments
WHERE status = $1
ORDER BY created_at DESC;
