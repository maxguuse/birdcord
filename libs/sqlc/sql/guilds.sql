-- name: CreateGuild :one
INSERT INTO guilds (
    discord_guild_id
) VALUES (
    $1
) RETURNING *;

-- name: CreateGuilds :execrows
INSERT INTO guilds ("discord_guild_id") 
VALUES (UNNEST(@discord_guild_ids::varchar[])) 
ON CONFLICT ("discord_guild_id") DO NOTHING;

-- name: GetGuildByID :one
SELECT * FROM guilds WHERE id = $1;

-- name: GetGuildByDiscordID :one
SELECT * FROM guilds WHERE discord_guild_id = $1;

-- name: DeleteGuildByID :exec
DELETE FROM guilds WHERE id = $1;

-- name: DeleteGuildByDiscordID :exec
DELETE FROM guilds WHERE discord_guild_id = $1;
