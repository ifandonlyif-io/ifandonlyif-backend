-- name: CreateUser :one
INSERT INTO users (
  full_name,
  wallet_address,
  country_code,
  email_address,
  twitter_name,
  image_uri,
  nonce
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByWalletAddress :one
SELECT * FROM users
WHERE wallet_address = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users;

-- name: UpdateUserEmailAddress :one
UPDATE users
SET email_address = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserTwitterName :one
UPDATE users
SET twitter_name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserKycDate :one
UPDATE users
SET kyc_date = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserNonce :one
UPDATE users
SET nonce = $2
WHERE wallet_address = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;