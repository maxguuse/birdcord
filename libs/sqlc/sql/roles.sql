-- name: CreateRole :one
INSERT INTO roles (guild_id, discord_role_id) VALUES ($1, $2) RETURNING *;

-- name: GetRoleByDiscordID :one
SELECT * FROM roles WHERE discord_role_id = $1;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1;

-- name: DeleteRoleByID :exec
DELETE FROM roles WHERE id = $1;

-- name: DeleteRoleByDiscordID :exec
DELETE FROM roles WHERE discord_role_id = $1;
