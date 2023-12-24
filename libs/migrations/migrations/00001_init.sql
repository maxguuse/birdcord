-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "guilds" (
    "id" serial PRIMARY KEY,
    "discord_guild_id" varchar(32) UNIQUE NOT NULL
);

CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "discord_user_id" varchar(32) UNIQUE NOT NULL
);

CREATE TABLE "messages" (
    "id" serial PRIMARY KEY,
    "discord_message_id" varchar(32) NOT NULL,
    "discord_channel_id" varchar(32) NOT NULL,
    CONSTRAINT unique_discord_message_and_channel UNIQUE ("discord_message_id", "discord_channel_id")
);

CREATE TABLE "polls" (
    "id" serial PRIMARY KEY,
    "title" varchar(100) NOT NULL,
    "is_active" bool DEFAULT true NOT NULL,
    "created_at" timestamp DEFAULT now(),
    "guild_id" int NOT NULL,
    "author_id" int
);

CREATE TABLE "poll_options" (
    "id" serial PRIMARY KEY,
    "title" varchar(100) NOT NULL,
    "poll_id" int NOT NULL
);

CREATE TABLE "poll_messages" (
    "id" serial PRIMARY KEY,
    "message_id" int NOT NULL,
    "poll_id" int NOT NULL,
    CONSTRAINT unique_message_and_poll UNIQUE ("message_id", "poll_id")
);

CREATE TABLE "poll_votes" (
    "id" serial PRIMARY KEY,
    "poll_id" int NOT NULL,
    "option_id" int NOT NULL,
    "user_id" int NOT NULL,
    CONSTRAINT unique_poll_and_user UNIQUE ("poll_id", "user_id")
);

ALTER TABLE "polls" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");
ALTER TABLE "polls" ADD FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id");
ALTER TABLE "poll_options" ADD FOREIGN KEY ("poll_id") REFERENCES "polls" ("id");
ALTER TABLE "poll_messages" ADD FOREIGN KEY ("poll_id") REFERENCES "polls" ("id");
ALTER TABLE "poll_messages" ADD FOREIGN KEY ("message_id") REFERENCES "messages" ("id") ON DELETE CASCADE;
ALTER TABLE "poll_votes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "poll_votes" ADD FOREIGN KEY ("poll_id") REFERENCES "polls" ("id");
ALTER TABLE "poll_votes" ADD FOREIGN KEY ("option_id") REFERENCES "poll_options" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "polls" DROP CONSTRAINT IF EXISTS "polls_author_id_fkey";
ALTER TABLE "polls" DROP CONSTRAINT IF EXISTS "polls_guild_id_fkey";
ALTER TABLE "poll_options" DROP CONSTRAINT IF EXISTS "poll_options_poll_id_fkey";
ALTER TABLE "poll_messages" DROP CONSTRAINT IF EXISTS "poll_messages_poll_id_fkey";
ALTER TABLE "poll_messages" DROP CONSTRAINT IF EXISTS "poll_messages_message_id_fkey";
ALTER TABLE "poll_votes" DROP CONSTRAINT IF EXISTS "poll_votes_user_id_fkey";
ALTER TABLE "poll_votes" DROP CONSTRAINT IF EXISTS "poll_votes_poll_id_fkey";
ALTER TABLE "poll_votes" DROP CONSTRAINT IF EXISTS "poll_votes_option_id_fkey";
DROP TABLE IF EXISTS "poll_votes";
DROP TABLE IF EXISTS "poll_messages";
DROP TABLE IF EXISTS "poll_options";
DROP TABLE IF EXISTS "polls";
DROP TABLE IF EXISTS "messages";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "guilds";
-- +goose StatementEnd
