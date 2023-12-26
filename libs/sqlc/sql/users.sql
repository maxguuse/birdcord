-- name: CreateUser :one
INSERT INTO users (
    discord_user_id
) VALUES (
    $1
) ON CONFLICT (discord_user_id) DO NOTHING RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByDiscordID :one
SELECT * FROM users WHERE discord_user_id = $1;

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUserByDiscordID :exec
DELETE FROM users WHERE discord_user_id = $1;