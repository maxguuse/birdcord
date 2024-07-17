-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE liveroles ADD CONSTRAINT unique_role_and_guild UNIQUE ("discord_role_id", "discord_guild_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE liveroles DROP CONSTRAINT unique_role_and_guild;
-- +goose StatementEnd
