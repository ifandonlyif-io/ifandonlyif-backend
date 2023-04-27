-- name: CreateUser :one
INSERT INTO users (
  full_name,
  wallet,
  country_code,
  email_address,
  twitter_name,
  image_uri,
  nonce
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUser :one
select ID, COALESCE(full_name),COALESCE(wallet),COALESCE(created_at),COALESCE(country_code),COALESCE(email_address),COALESCE(kyc_date),COALESCE(twitter_name),COALESCE(blockpass_id),COALESCE(image_uri),COALESCE(nonce)
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByWalletAddress :one
select ID,COALESCE(full_name),COALESCE(wallet),COALESCE(nonce)
FROM users
WHERE wallet = $1 LIMIT 1;

-- name: GetUserForUpdate :one
select *
FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
select *
FROM users;

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
WHERE wallet = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;