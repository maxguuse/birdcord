-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "polls" (
    "id" serial PRIMARY KEY,
    "title" varchar,
    "discord_token" varchar UNIQUE
);

CREATE TABLE "polls_options" (
    "id" serial PRIMARY KEY,
    "title" varchar,
    "poll_id" int
);

CREATE TABLE "voted_users" (
    "id" serial PRIMARY KEY,
    "discord_id" varchar,
    "option_id" int,
    "poll_id" int
);

ALTER TABLE "polls_options" ADD FOREIGN KEY ("poll_id") REFERENCES "polls" ("id");
ALTER TABLE "voted_users" ADD FOREIGN KEY ("poll_id") REFERENCES "polls" ("id");
ALTER TABLE "voted_users" ADD FOREIGN KEY ("option_id") REFERENCES "polls_options" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS "voted_users";
DROP TABLE IF EXISTS "polls_options";
DROP TABLE IF EXISTS "polls";
-- +goose StatementEnd
