-- name: CreatePollOption :one
INSERT INTO poll_options (
    title, 
    poll_id
) VALUES (
    $1, $2
) RETURNING *;