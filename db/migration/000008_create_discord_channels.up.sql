DROP TABLE IF EXISTS discord_channels;

CREATE TABLE discord_channels (
    "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar NOT NULL,
    "guild_id" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

DROP TABLE IF EXISTS appliances_from_discord_channel;

CREATE TABLE appliances_from_discord_channel (
    "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    "channel_name" varchar NOT NULL,
    "guild_id" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "verified_at" timestamptz,
    "is_approved" bool
)