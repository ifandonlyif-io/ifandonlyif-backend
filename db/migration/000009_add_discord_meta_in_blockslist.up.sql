ALTER TABLE report_blocklists
    ADD COLUMN "guild_id" varchar NOT NULL DEFAULT '',
    ADD COLUMN "guild_name" varchar NOT NULL DEFAULT '',
    ADD COLUMN "reporter_name" varchar NOT NULL DEFAULT '',
    ADD COLUMN "reporter_avatar" varchar NOT NULL DEFAULT '';