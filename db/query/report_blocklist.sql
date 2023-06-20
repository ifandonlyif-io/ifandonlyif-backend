-- name: CreateReportBlocklist :one
INSERT INTO report_blocklists (
  http_address, user_wallet_address, guild_id, guild_name, reporter_name, reporter_avatar
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetReportBlocklist :one
SELECT * FROM report_blocklists
WHERE id = $1 LIMIT 1;

-- name: GetReportBlocklistByUrl :many
SELECT * FROM report_blocklists
WHERE http_address = $1;

-- name: GetReportBlocklistUpdate :one
SELECT * FROM report_blocklists
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListReportBlocklists :many
SELECT * FROM report_blocklists;

-- name: VerifyBlocklist :one
UPDATE report_blocklists
SET verified_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DisproveBlocklist :one
UPDATE report_blocklists
SET disproved_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReportBlocklist :exec
DELETE FROM report_blocklists
WHERE id = $1;

-- name: GetBlocklistByUri :one
SELECT * FROM report_blocklists 
WHERE http_address = $1;

-- name: ListVerifiedBlocklists :many
SELECT * FROM report_blocklists
WHERE verified_at is NOT NULL;

-- name: ListDisprovedBlocklists :many
SELECT * FROM report_blocklists
WHERE disproved_at is NOT NULL;

-- name: ListUnreviewedBlocklists :many
SELECT * FROM report_blocklists
WHERE disproved_at is NULL
AND verified_at is NULL;

-- name: CheckExistBlocklists :one
SELECT id FROM report_blocklists WHERE http_address = $1;

-- name: CheckBlocklists :one
SELECT 
    CASE 
        WHEN verified_at IS NOT NULL THEN 'verified_at' 
        WHEN disproved_at IS NOT NULL THEN 'disproved_at'
    END AS review_status,
    COALESCE(verified_at, disproved_at) AS time_stamp
FROM report_blocklists
WHERE id = $1;