-- name: GetGuildByDiscordID :one
SELECT * FROM guilds WHERE discord_guild_id = $1;

-- name: CreateGuilds :execrows
INSERT INTO guilds ("discord_guild_id") 
VALUES (UNNEST(@discord_guild_ids::varchar[])) 
ON CONFLICT ("discord_guild_id") DO NOTHING;
