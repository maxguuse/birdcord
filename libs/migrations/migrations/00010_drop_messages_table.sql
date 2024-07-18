-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE poll_messages ADD COLUMN discord_message_id bigint;
ALTER TABLE poll_messages ADD COLUMN discord_channel_id bigint;

UPDATE poll_messages
SET discord_message_id = subquery.discord_message_id,
    discord_channel_id = subquery.discord_channel_id
FROM (
    SELECT id, discord_message_id::bigint AS discord_message_id, discord_channel_id::bigint AS discord_channel_id
    FROM messages
) AS subquery
WHERE subquery.id = poll_messages.message_id;

ALTER TABLE poll_messages DROP COLUMN message_id;
DROP TABLE IF EXISTS "messages";
-- +goose StatementEnd

-- poll_messages
-- message_id

-- messages
-- discord_message_id
-- discord_channel_id

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

CREATE TABLE "messages" (
    "id" serial PRIMARY KEY,
    "discord_message_id" varchar(32) NOT NULL,
    "discord_channel_id" varchar(32) NOT NULL,
    CONSTRAINT unique_discord_message_and_channel UNIQUE ("discord_message_id", "discord_channel_id")
);

ALTER TABLE poll_messages ADD COLUMN message_id int;

WITH inserted_messages AS (
    INSERT INTO messages (discord_message_id, discord_channel_id)
    SELECT DISTINCT poll_messages.discord_message_id::varchar(32), poll_messages.discord_channel_id::varchar(32)
    FROM poll_messages
    RETURNING id
)

UPDATE poll_messages
SET message_id = inserted_messages.id
FROM inserted_messages;

ALTER TABLE poll_messages DROP COLUMN discord_message_id;
ALTER TABLE poll_messages DROP COLUMN discord_channel_id;

ALTER TABLE "poll_messages" ADD FOREIGN KEY ("message_id") REFERENCES "messages" ("id") ON DELETE CASCADE;
-- +goose StatementEnd
