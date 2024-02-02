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

CREATE TABLE "given_liveroles" (
    "id" serial PRIMARY KEY,
    "guild_id" int NOT NULL,
    "user_id" int NOT NULL,
    "liverole_id" int NOT NULL
);

ALTER TABLE "liveroles" ADD FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id");
ALTER TABLE "liveroles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");
ALTER TABLE "given_liveroles" ADD FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id");
ALTER TABLE "given_liveroles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "given_liveroles" ADD FOREIGN KEY ("liverole_id") REFERENCES "liveroles" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "liveroles" DROP CONSTRAINT IF EXISTS "liveroles_guild_id_fkey";
ALTER TABLE "liveroles" DROP CONSTRAINT IF EXISTS "liveroles_role_id_fkey";
ALTER TABLE "given_liveroles" DROP CONSTRAINT IF EXISTS "given_liveroles_guild_id_fkey";
ALTER TABLE "given_liveroles" DROP CONSTRAINT IF EXISTS "given_liveroles_user_id_fkey";
ALTER TABLE "given_liveroles" DROP CONSTRAINT IF EXISTS "given_liveroles_liverole_id_fkey";

DROP TABLE IF EXISTS "given_liveroles";
DROP TABLE IF EXISTS "liveroles";
DROP TABLE IF EXISTS "roles";
-- +goose StatementEnd
