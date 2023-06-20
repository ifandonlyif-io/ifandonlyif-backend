// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	CheckBlocklists(ctx context.Context, id uuid.UUID) (CheckBlocklistsRow, error)
	CheckExistBlocklists(ctx context.Context, httpAddress string) (uuid.UUID, error)
	CreateAppliance(ctx context.Context, arg CreateApplianceParams) (AppliancesFromDiscordChannel, error)
	CreateChannel(ctx context.Context, arg CreateChannelParams) (DiscordChannel, error)
	CreateGasPrice(ctx context.Context, average sql.NullInt32) (GasPrice, error)
	CreateIffNft(ctx context.Context, arg CreateIffNftParams) (IffNft, error)
	CreateReportBlocklist(ctx context.Context, arg CreateReportBlocklistParams) (ReportBlocklist, error)
	CreateReportWhitelist(ctx context.Context, arg CreateReportWhitelistParams) (ReportWhitelist, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteChannelById(ctx context.Context, id uuid.UUID) error
	DeleteReportBlocklist(ctx context.Context, id uuid.UUID) error
	DeleteReportWhitelist(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DisproveBlocklist(ctx context.Context, id uuid.UUID) (ReportBlocklist, error)
	GetAllAppliances(ctx context.Context) ([]AppliancesFromDiscordChannel, error)
	GetAllChannels(ctx context.Context) ([]DiscordChannel, error)
	GetApplianceByGuildId(ctx context.Context, guildID string) (AppliancesFromDiscordChannel, error)
	// :param id: string
	GetApplianceChannelById(ctx context.Context, id uuid.UUID) (AppliancesFromDiscordChannel, error)
	GetAveragePriceByLastDay(ctx context.Context) ([]GetAveragePriceByLastDayRow, error)
	GetBlocklistByUri(ctx context.Context, httpAddress string) (ReportBlocklist, error)
	GetChannelById(ctx context.Context, id uuid.UUID) (DiscordChannel, error)
	GetChannelsByGuildId(ctx context.Context, guildID string) (DiscordChannel, error)
	GetIffNftForUpdate(ctx context.Context, id uuid.UUID) (IffNft, error)
	GetIffNfts(ctx context.Context, id uuid.UUID) (IffNft, error)
	GetReportBlocklist(ctx context.Context, id uuid.UUID) (ReportBlocklist, error)
	GetReportBlocklistByUrl(ctx context.Context, httpAddress string) ([]ReportBlocklist, error)
	GetReportBlocklistUpdate(ctx context.Context, id uuid.UUID) (ReportBlocklist, error)
	GetReportWhitelis(ctx context.Context, id uuid.UUID) (ReportWhitelist, error)
	GetReportWhitelistUpdate(ctx context.Context, id uuid.UUID) (ReportWhitelist, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByWalletAddress(ctx context.Context, wallet sql.NullString) (GetUserByWalletAddressRow, error)
	GetUserForUpdate(ctx context.Context, id uuid.UUID) (User, error)
	ListDisprovedBlocklists(ctx context.Context) ([]ReportBlocklist, error)
	ListIffNfts(ctx context.Context) ([]IffNft, error)
	ListNftProjects(ctx context.Context) ([]NftProject, error)
	ListReportBlocklists(ctx context.Context) ([]ReportBlocklist, error)
	ListReportWhitelist(ctx context.Context) ([]ReportWhitelist, error)
	ListUnreviewedBlocklists(ctx context.Context) ([]ReportBlocklist, error)
	ListUsers(ctx context.Context) ([]User, error)
	ListVerifiedBlocklists(ctx context.Context) ([]ReportBlocklist, error)
	// :param id: string
	// :param isApproved: bool
	UpdateApplianceChannel(ctx context.Context, arg UpdateApplianceChannelParams) (AppliancesFromDiscordChannel, error)
	UpdateReportWhitelisVerified(ctx context.Context, arg UpdateReportWhitelisVerifiedParams) (ReportWhitelist, error)
	UpdateUserEmailAddress(ctx context.Context, arg UpdateUserEmailAddressParams) (User, error)
	UpdateUserKycDate(ctx context.Context, arg UpdateUserKycDateParams) (User, error)
	UpdateUserNonce(ctx context.Context, arg UpdateUserNonceParams) (User, error)
	UpdateUserTwitterName(ctx context.Context, arg UpdateUserTwitterNameParams) (User, error)
	VerifyBlocklist(ctx context.Context, id uuid.UUID) (ReportBlocklist, error)
}

var _ Querier = (*Queries)(nil)
