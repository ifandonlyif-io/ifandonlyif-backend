// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: discord_channel.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createAppliance = `-- name: CreateAppliance :one
INSERT INTO appliances_from_discord_channel ("channel_name", "guild_id")
VALUES ($1, $2) returning id, channel_name, guild_id, created_at, verified_at, is_approved
`

type CreateApplianceParams struct {
	ChannelName string `json:"channelName"`
	GuildID     string `json:"guildID"`
}

func (q *Queries) CreateAppliance(ctx context.Context, arg CreateApplianceParams) (AppliancesFromDiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, createAppliance, arg.ChannelName, arg.GuildID)
	var i AppliancesFromDiscordChannel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.GuildID,
		&i.CreatedAt,
		&i.VerifiedAt,
		&i.IsApproved,
	)
	return i, err
}

const createChannel = `-- name: CreateChannel :one
INSERT INTO discord_channels (
  name, guild_id
) values (
  $1, $2
) RETURNING id, name, guild_id, created_at
`

type CreateChannelParams struct {
	Name    string `json:"name"`
	GuildID string `json:"guildID"`
}

func (q *Queries) CreateChannel(ctx context.Context, arg CreateChannelParams) (DiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, createChannel, arg.Name, arg.GuildID)
	var i DiscordChannel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.GuildID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteChannelById = `-- name: DeleteChannelById :exec
DELETE FROM discord_channels WHERE id = $1
`

func (q *Queries) DeleteChannelById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteChannelById, id)
	return err
}

const getAllAppliances = `-- name: GetAllAppliances :many
SELECT id, channel_name, guild_id, created_at, verified_at, is_approved FROM appliances_from_discord_channel ORDER BY "created_at" DESC
`

func (q *Queries) GetAllAppliances(ctx context.Context) ([]AppliancesFromDiscordChannel, error) {
	rows, err := q.db.QueryContext(ctx, getAllAppliances)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AppliancesFromDiscordChannel{}
	for rows.Next() {
		var i AppliancesFromDiscordChannel
		if err := rows.Scan(
			&i.ID,
			&i.ChannelName,
			&i.GuildID,
			&i.CreatedAt,
			&i.VerifiedAt,
			&i.IsApproved,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllChannels = `-- name: GetAllChannels :many
SELECT id, name, guild_id, created_at FROM discord_channels
`

func (q *Queries) GetAllChannels(ctx context.Context) ([]DiscordChannel, error) {
	rows, err := q.db.QueryContext(ctx, getAllChannels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []DiscordChannel{}
	for rows.Next() {
		var i DiscordChannel
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.GuildID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApplianceByGuildId = `-- name: GetApplianceByGuildId :one
SELECT id, channel_name, guild_id, created_at, verified_at, is_approved FROM appliances_from_discord_channel WHERE guild_id = $1
`

func (q *Queries) GetApplianceByGuildId(ctx context.Context, guildID string) (AppliancesFromDiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, getApplianceByGuildId, guildID)
	var i AppliancesFromDiscordChannel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.GuildID,
		&i.CreatedAt,
		&i.VerifiedAt,
		&i.IsApproved,
	)
	return i, err
}

const getApplianceChannelById = `-- name: GetApplianceChannelById :one
SELECT id, channel_name, guild_id, created_at, verified_at, is_approved FROM appliances_from_discord_channel WHERE "id" = $1
`

// :param id: string
func (q *Queries) GetApplianceChannelById(ctx context.Context, id uuid.UUID) (AppliancesFromDiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, getApplianceChannelById, id)
	var i AppliancesFromDiscordChannel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.GuildID,
		&i.CreatedAt,
		&i.VerifiedAt,
		&i.IsApproved,
	)
	return i, err
}

const getChannelById = `-- name: GetChannelById :one
SELECT id, name, guild_id, created_at FROM discord_channels
WHERE id = $1
`

func (q *Queries) GetChannelById(ctx context.Context, id uuid.UUID) (DiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, getChannelById, id)
	var i DiscordChannel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.GuildID,
		&i.CreatedAt,
	)
	return i, err
}

const getChannelsByGuildId = `-- name: GetChannelsByGuildId :one
SELECT id, name, guild_id, created_at FROM discord_channels
WHERE guild_id = $1
`

func (q *Queries) GetChannelsByGuildId(ctx context.Context, guildID string) (DiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, getChannelsByGuildId, guildID)
	var i DiscordChannel
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.GuildID,
		&i.CreatedAt,
	)
	return i, err
}

const updateApplianceChannel = `-- name: UpdateApplianceChannel :one
UPDATE appliances_from_discord_channel
SET "is_approved" = $2, "verified_at" = now()
WHERE "id" = $1
RETURNING id, channel_name, guild_id, created_at, verified_at, is_approved
`

type UpdateApplianceChannelParams struct {
	ID         uuid.UUID    `json:"id"`
	IsApproved sql.NullBool `json:"isApproved"`
}

// :param id: string
// :param isApproved: bool
func (q *Queries) UpdateApplianceChannel(ctx context.Context, arg UpdateApplianceChannelParams) (AppliancesFromDiscordChannel, error) {
	row := q.db.QueryRowContext(ctx, updateApplianceChannel, arg.ID, arg.IsApproved)
	var i AppliancesFromDiscordChannel
	err := row.Scan(
		&i.ID,
		&i.ChannelName,
		&i.GuildID,
		&i.CreatedAt,
		&i.VerifiedAt,
		&i.IsApproved,
	)
	return i, err
}