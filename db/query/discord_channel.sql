-- name: CreateChannel :one
INSERT INTO discord_channels (
  name, guild_id
) values (
  $1, $2
) RETURNING *;

-- name: GetAllChannels :many
SELECT * FROM discord_channels;

-- name: GetChannelById :one
SELECT * FROM discord_channels
WHERE id = $1;

-- name: GetChannelsByGuildId :one
SELECT * FROM discord_channels
WHERE guild_id = $1;

-- name: DeleteChannelById :exec
DELETE FROM discord_channels WHERE id = $1;

-- name: CreateAppliance :one
INSERT INTO appliances_from_discord_channel ("channel_name", "guild_id")
VALUES ($1, $2) returning *;

-- name: GetAllAppliances :many
SELECT * FROM appliances_from_discord_channel ORDER BY "created_at" DESC;

-- name: GetApplianceChannelById :one
-- :param id: string
SELECT * FROM appliances_from_discord_channel WHERE "id" = $1;

-- name: GetApplianceByGuildId :one
SELECT * FROM appliances_from_discord_channel WHERE guild_id = $1;

-- name: UpdateApplianceChannel :one
-- :param id: string
-- :param isApproved: bool
UPDATE appliances_from_discord_channel
SET "is_approved" = $2, "verified_at" = now()
WHERE "id" = $1
RETURNING *;

-- name: LockDiscordChannel :one
-- :param id: string
UPDATE discord_channels
SET "locked_at" = now()
WHERE "id" = $1
RETURNING *;

-- name: UnlockDiscordChannel :one
-- :param id: string
UPDATE discord_channels
SET "locked_at" = null
WHERE "id" = $1
RETURNING *;