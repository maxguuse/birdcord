-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE liveroles ADD COLUMN discord_role_id bigint;
ALTER TABLE liveroles ADD COLUMN discord_guild_id bigint;

UPDATE liveroles
SET discord_role_id = subquery.discord_role_id,
    discord_guild_id = subquery.guild_id
FROM (
    SELECT id, discord_role_id::bigint AS discord_role_id, guild_id::bigint AS guild_id
    FROM roles
) AS subquery
WHERE subquery.id = liveroles.role_id;

ALTER TABLE liveroles DROP COLUMN role_id;
DROP TABLE IF EXISTS "roles";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

CREATE TABLE "roles" (
    "id" serial PRIMARY KEY,
    "guild_id" bigint NOT NULL,
    "discord_role_id" varchar(32) UNIQUE NOT NULL
);

ALTER TABLE liveroles ADD COLUMN role_id int;

WITH inserted_roles AS (
    INSERT INTO roles (discord_role_id, guild_id)
    SELECT DISTINCT liveroles.discord_role_id::varchar(32), liveroles.discord_guild_id
    FROM liveroles
    RETURNING id
)

UPDATE liveroles
SET role_id = inserted_roles.id
FROM inserted_roles;

ALTER TABLE liveroles DROP COLUMN discord_role_id;
ALTER TABLE liveroles DROP COLUMN discord_guild_id;

ALTER TABLE "liveroles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
-- +goose StatementEnd
