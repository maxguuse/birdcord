-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "polls" ADD COLUMN "discord_author_id" varchar;
ALTER TABLE "polls" ADD COLUMN "discord_guild_id" varchar;
ALTER TABLE "polls" ADD COLUMN "active" bool DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "polls" DROP COLUMN "discord_author_id";
ALTER TABLE "polls" DROP COLUMN "discord_guild_id";
ALTER TABLE "polls" DROP COLUMN "active";
-- +goose StatementEnd
