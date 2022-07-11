-- name: CreateIffNft :one
INSERT INTO iff_nfts (
  project_id,
  user_wallet_address,
  nft_projects_contract_address,
  nft_projects_collection_name,
  mint_date,
  mint_transaction
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetIffNfts :one
SELECT * FROM iff_nfts
WHERE id = $1 LIMIT 1;

-- name: GetIffNftForUpdate :one
SELECT * FROM iff_nfts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListIffNfts :many
SELECT * FROM iff_nfts;
