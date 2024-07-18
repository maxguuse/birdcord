-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE IF EXISTS "guilds";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

CREATE TABLE "guilds" (
    "id" serial PRIMARY KEY,
    "discord_guild_id" varchar(32) UNIQUE NOT NULL
);
-- +goose StatementEnd
