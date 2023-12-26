-- name: CreatePollMessage :one
INSERT INTO poll_messages (
    poll_id, 
    message_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetFullPollMessages :many
SELECT pm.id, pm.message_id, m.discord_message_id, m.discord_channel_id
FROM poll_messages pm 
LEFT JOIN messages m ON pm.message_id = m.id
WHERE pm.poll_id = $1;

-- name: GetPollMessages :many
SELECT * FROM poll_messages WHERE poll_id = $1;

-- name: GetPollMessageByMessageId :one
SELECT * FROM poll_messages WHERE message_id = $1;

-- name: DeletePollMessageById :exec
DELETE FROM poll_messages WHERE id = $1;