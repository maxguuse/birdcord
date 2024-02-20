-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "roles" (
    "id" serial PRIMARY KEY,
    "discord_role_id" varchar(32) UNIQUE NOT NULL
);

CREATE TABLE "liveroles" (
    "id" serial PRIMARY KEY,
    "guild_id" int NOT NULL,
    "role_id" int NOT NULL,
    CONSTRAINT unique_liverole_for_guild UNIQUE ("guild_id", "role_id")
);

ALTER TABLE "liveroles" ADD FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id");
ALTER TABLE "liveroles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "liveroles" DROP CONSTRAINT IF EXISTS "liveroles_guild_id_fkey";
ALTER TABLE "liveroles" DROP CONSTRAINT IF EXISTS "liveroles_role_id_fkey";

DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "liveroles";
-- +goose StatementEnd
