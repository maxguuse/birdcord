-- name: GetLiveRolesByGuildID :many
SELECT * FROM liveroles 
LEFT JOIN roles ON liveroles.role_id = roles.id
WHERE roles.guild_id = $1;

-- name: GetLiveRoleByDiscordID :one
SELECT * FROM liveroles
LEFT JOIN roles ON liveroles.role_id = roles.id
WHERE roles.discord_role_id = $1
AND roles.guild_id = $2;

-- name: CreateLiveRole :one
INSERT INTO liveroles (role_id) VALUES ($1) RETURNING *;

-- name: DeleteLiveRoleByRoleID :exec
DELETE FROM liveroles WHERE role_id = $1;