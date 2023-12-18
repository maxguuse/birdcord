-- name: GetUserByDiscordID :one
SELECT * FROM users WHERE discord_user_id = $1;

-- name: CreateUser :one
INSERT INTO users (
    discord_user_id
) VALUES (
    $1
) ON CONFLICT (discord_user_id) DO NOTHING RETURNING *;