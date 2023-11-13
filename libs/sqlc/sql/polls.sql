-- name: CreatePoll :one
INSERT INTO polls (
    title, discord_id, discord_author_id, discord_guild_id, channel_id
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetToken :one
SELECT discord_id FROM polls
                     WHERE id = $1;

-- name: GetPoll :one
SELECT * FROM polls
         WHERE id = $1;

-- name: GetActivePolls :many
SELECT id, title FROM polls
         WHERE discord_guild_id = $1 AND active = true;

-- name: StopPoll :exec
UPDATE polls
    SET active = FALSE
    WHERE id = $1;

-- name: RemoveVotedUsers :exec
DELETE FROM voted_users WHERE poll_id = $1;

-- name: RemovePollOptions :exec
DELETE FROM polls_options WHERE poll_id = $1;

-- name: RemovePoll :exec
DELETE FROM polls WHERE id = $1;
