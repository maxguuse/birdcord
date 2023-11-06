-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE polls RENAME COLUMN discord_token TO discord_id;
ALTER TABLE polls ADD COLUMN channel_id varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE polls RENAME COLUMN discord_id TO discord_token;
ALTER TABLE polls DROP COLUMN channel_id;
-- +goose StatementEnd
