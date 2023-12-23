-- name: CreatePollOption :one
INSERT INTO poll_options (
    title, 
    poll_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: CreatePollOptions :many
INSERT INTO poll_options ("title", "poll_id") 
VALUES (UNNEST(@titles::varchar[]), @poll_id)
RETURNING sqlc.embed(poll_options);

-- name: GetPollOptions :many
SELECT * FROM poll_options
WHERE poll_id = $1;