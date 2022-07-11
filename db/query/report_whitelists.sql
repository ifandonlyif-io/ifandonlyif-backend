-- name: CreateReportWhitelist :one
INSERT INTO report_whitelists (
  http_address,
  user_wallet_address
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetReportWhitelis :one
SELECT * FROM report_whitelists
WHERE id = $1 LIMIT 1;

-- name: GetReportWhitelistUpdate :one
SELECT * FROM report_whitelists
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListReportWhitelist :many
SELECT * FROM report_whitelists;

-- name: UpdateReportWhitelisVerified :one
UPDATE report_whitelists
SET verified_at = $2
WHERE id = $1
RETURNING *;

-- name: DeleteReportWhitelist :exec
DELETE FROM report_whitelists
WHERE id = $1;