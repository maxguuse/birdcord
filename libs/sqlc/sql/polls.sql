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