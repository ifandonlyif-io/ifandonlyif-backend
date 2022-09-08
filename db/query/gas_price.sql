-- name: CreateGasPrice :one
INSERT INTO gas_prices (
    average
) VALUES (
  $1
) RETURNING *;

-- name: GetAveragePriceByLastDay :many
  SELECT created_at,
    average
  FROM gas_prices
  ORDER BY created_at DESC
  LIMIT 24;