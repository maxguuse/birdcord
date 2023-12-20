-- name: CreateMessage :one
INSERT INTO messages (
    discord_message_id,
    discord_channel_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetMessageById :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessageByDiscordID :one
SELECT * FROM messages WHERE discord_message_id = $1;

-- name: DeleteMessageById :exec
DELETE FROM messages WHERE id = $1;