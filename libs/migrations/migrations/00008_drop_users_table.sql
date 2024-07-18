-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP TABLE IF EXISTS "users";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "discord_user_id" varchar(32) UNIQUE NOT NULL
);
-- +goose StatementEnd
