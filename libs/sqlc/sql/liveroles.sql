-- name: GetLiveRolesByGuildID :many
SELECT * FROM liveroles WHERE guild_id = $1;

-- name: CreateLiveRole :one
INSERT INTO liveroles (guild_id, role_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteLiveRoleByRoleID :exec
DELETE FROM liveroles WHERE role_id = $1;