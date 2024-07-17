-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE poll_votes ADD CONSTRAINT unique_poll_and_user UNIQUE ("poll_id", "user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE poll_votes DROP CONSTRAINT unique_poll_and_user;
-- +goose StatementEnd
