-- name: CreatePollOption :one
INSERT INTO poll_options (
    title, 
    poll_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetPollOptions :many
SELECT * FROM poll_options
WHERE poll_id = $1;