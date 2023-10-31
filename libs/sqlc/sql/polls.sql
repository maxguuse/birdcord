-- name: CreatePoll :one
INSERT INTO polls (
    title,
    discord_token
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetToken :one
SELECT discord_token FROM polls
                     WHERE id = $1;

-- name: GetPoll :one
SELECT * FROM polls
         WHERE id = $1;

-- name: GetPolls :many
SELECT id, title FROM polls;