// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type IffNft struct {
	ID                         uuid.UUID `json:"id"`
	ProjectID                  int64     `json:"project_id"`
	UserWalletAddress          string    `json:"user_wallet_address"`
	NftProjectsContractAddress string    `json:"nft_projects_contract_address"`
	NftProjectsCollectionName  string    `json:"nft_projects_collection_name"`
	MintDate                   time.Time `json:"mint_date"`
	MintTransaction            string    `json:"mint_transaction"`
}

type NftProject struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ContractAddress string    `json:"contract_address"`
	CollectionName  string    `json:"collection_name"`
}

type ReportBlocklist struct {
	ID                uuid.UUID      `json:"id"`
	HttpAddress       string         `json:"http_address"`
	VerifiedAt        sql.NullTime   `json:"verified_at"`
	UserWalletAddress sql.NullString `json:"user_wallet_address"`
	CreatedAt         time.Time      `json:"created_at"`
}

type ReportWhitelist struct {
	ID                uuid.UUID      `json:"id"`
	HttpAddress       string         `json:"http_address"`
	VerifiedAt        sql.NullTime   `json:"verified_at"`
	UserWalletAddress sql.NullString `json:"user_wallet_address"`
	CreatedAt         time.Time      `json:"created_at"`
}

type User struct {
	ID            uuid.UUID      `json:"id"`
	FullName      sql.NullString `json:"full_name"`
	WalletAddress sql.NullString `json:"wallet_address"`
	CreatedAt     sql.NullTime   `json:"created_at"`
	CountryCode   sql.NullString `json:"country_code"`
	EmailAddress  sql.NullString `json:"email_address"`
	KycDate       sql.NullTime   `json:"kyc_date"`
	TwitterName   sql.NullString `json:"twitter_name"`
	BlockpassID   sql.NullInt64  `json:"blockpass_id"`
	ImageUri      sql.NullString `json:"image_uri"`
	Nonce         sql.NullString `json:"nonce"`
}
