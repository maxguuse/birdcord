-- name: CreateMessage :one
INSERT INTO messages (
    discord_message_id,
    discord_channel_id
) VALUES (
    $1, $2
) RETURNING *;