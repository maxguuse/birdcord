-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE polls ADD COLUMN discord_guild_id bigint;

UPDATE polls
SET discord_guild_id = (
    SELECT discord_guild_id::bigint
    FROM guilds
    WHERE guilds.id = polls.guild_id
);

ALTER TABLE polls DROP COLUMN guild_id;

ALTER TABLE polls RENAME COLUMN discord_guild_id TO guild_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE polls ADD COLUMN internal_guild_id int;

WITH inserted_guilds AS (
    INSERT INTO guilds (discord_guild_id)
    SELECT DISTINCT polls.guild_id::varchar(32)
    FROM polls
    LEFT JOIN guilds ON guilds.discord_guild_id::bigint = polls.guild_id
    WHERE guilds.id IS NULL
    RETURNING id, discord_guild_id
), all_guilds AS (
    SELECT id, discord_guild_id FROM guilds
    UNION ALL
    SELECT id, discord_guild_id FROM inserted_guilds
)

UPDATE polls
SET internal_guild_id = (
    SELECT id
    FROM all_guilds
    WHERE all_guilds.discord_guild_id::bigint = polls.guild_id
);

ALTER TABLE polls DROP COLUMN guild_id;

ALTER TABLE polls RENAME COLUMN internal_guild_id to guild_id;

ALTER TABLE "polls" ADD FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id");
-- +goose StatementEnd
