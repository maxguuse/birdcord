-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE polls ADD COLUMN discord_author_id bigint;
ALTER TABLE poll_votes ADD COLUMN discord_user_id bigint;

UPDATE polls
SET discord_author_id = (
    SELECT discord_user_id::bigint
    FROM users
    WHERE users.id = polls.author_id
);

UPDATE poll_votes
SET discord_user_id = (
    SELECT discord_user_id::bigint
    FROM users
    WHERE users.id = poll_votes.user_id
);

ALTER TABLE polls DROP COLUMN author_id;
ALTER TABLE poll_votes DROP COLUMN user_id;

ALTER TABLE polls RENAME COLUMN discord_author_id TO author_id;
ALTER TABLE poll_votes RENAME COLUMN discord_user_id TO user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE polls ADD COLUMN internal_user_id int;
ALTER TABLE poll_votes ADD COLUMN internal_user_id int;

WITH inserted_users AS (
    INSERT INTO users (discord_user_id)
    (
        SELECT DISTINCT polls.author_id::varchar(32)
        FROM polls
        LEFT JOIN users ON users.discord_user_id::bigint = polls.author_id
        WHERE users.id IS NULL
        UNION
        SELECT DISTINCT poll_votes.user_id::varchar(32)
        FROM poll_votes
        LEFT JOIN users ON users.discord_user_id::bigint = poll_votes.user_id
        WHERE users.id IS NULL
    )
    RETURNING id, discord_user_id
), all_users AS (
    SELECT id, discord_user_id FROM users
    UNION ALL
    SELECT id, discord_user_id FROM inserted_users
)

UPDATE polls
SET internal_user_id = (
    SELECT id
    FROM all_users
    WHERE all_users.discord_user_id::bigint = polls.author_id
);

UPDATE poll_votes
SET internal_user_id = (
    SELECT id
    FROM users
    WHERE users.discord_user_id::bigint = poll_votes.user_id
);

ALTER TABLE polls DROP COLUMN author_id;
ALTER TABLE poll_votes DROP COLUMN user_id;

ALTER TABLE polls RENAME COLUMN internal_user_id to author_id;
ALTER TABLE poll_votes RENAME COLUMN internal_user_id to user_id;

ALTER TABLE "polls" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");
ALTER TABLE "poll_votes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
-- +goose StatementEnd
