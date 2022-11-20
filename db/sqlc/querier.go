// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	CreateGasPrice(ctx context.Context, average sql.NullInt32) (GasPrice, error)
	CreateIffNft(ctx context.Context, arg CreateIffNftParams) (IffNft, error)
	CreateReportBlocklist(ctx context.Context, arg CreateReportBlocklistParams) (ReportBlocklist, error)
	CreateReportWhitelist(ctx context.Context, arg CreateReportWhitelistParams) (ReportWhitelist, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteReportBlocklist(ctx context.Context, id uuid.UUID) error
	DeleteReportWhitelist(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAveragePriceByLastDay(ctx context.Context) ([]GetAveragePriceByLastDayRow, error)
	GetBlocklistByUri(ctx context.Context, httpAddress string) (ReportBlocklist, error)
	GetIffNftForUpdate(ctx context.Context, id uuid.UUID) (IffNft, error)
	GetIffNfts(ctx context.Context, id uuid.UUID) (IffNft, error)
	GetReportBlocklist(ctx context.Context, id uuid.UUID) (ReportBlocklist, error)
	GetReportBlocklistUpdate(ctx context.Context, id uuid.UUID) (ReportBlocklist, error)
	GetReportWhitelis(ctx context.Context, id uuid.UUID) (ReportWhitelist, error)
	GetReportWhitelistUpdate(ctx context.Context, id uuid.UUID) (ReportWhitelist, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id uuid.UUID) (GetUserRow, error)
	GetUserByWalletAddress(ctx context.Context, wallet sql.NullString) (GetUserByWalletAddressRow, error)
	GetUserForUpdate(ctx context.Context, id uuid.UUID) (User, error)
	ListIffNfts(ctx context.Context) ([]IffNft, error)
	ListReportBlocklists(ctx context.Context) ([]ReportBlocklist, error)
	ListReportWhitelist(ctx context.Context) ([]ReportWhitelist, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateReportBlocklistVerified(ctx context.Context, arg UpdateReportBlocklistVerifiedParams) (ReportBlocklist, error)
	UpdateReportWhitelisVerified(ctx context.Context, arg UpdateReportWhitelisVerifiedParams) (ReportWhitelist, error)
	UpdateUserEmailAddress(ctx context.Context, arg UpdateUserEmailAddressParams) (User, error)
	UpdateUserKycDate(ctx context.Context, arg UpdateUserKycDateParams) (User, error)
	UpdateUserNonce(ctx context.Context, arg UpdateUserNonceParams) (User, error)
	UpdateUserTwitterName(ctx context.Context, arg UpdateUserTwitterNameParams) (User, error)
}

var _ Querier = (*Queries)(nil)
