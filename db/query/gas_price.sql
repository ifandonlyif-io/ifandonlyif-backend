-- name: CreateGasPrice :one
INSERT INTO gas_prices (
    average
) VALUES (
  $1
) RETURNING *;

-- name: GetAveragePriceByLastDay :many
 SELECT COALESCE(average) AS average,
  COALESCE(created_at) AS created_at
  FROM gas_prices
  ORDER BY created_at DESC
  LIMIT 24;