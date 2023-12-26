-- name: CreatePoll :one
INSERT INTO polls (
    title, 
    author_id, 
    guild_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPoll :one
SELECT * FROM polls WHERE id = $1;

-- name: GetActivePolls :many
SELECT * FROM polls WHERE is_active = true AND guild_id = $1 AND author_id = $2;

-- name: UpdatePollStatus :exec
UPDATE polls SET "is_active" = $2 WHERE id = $1;