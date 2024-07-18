-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "liveroles" ADD CONSTRAINT unique_role_id_per_liverole UNIQUE (role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE "liveroles" DROP CONSTRAINT IF EXISTS unique_role_id_per_liverole;
-- +goose StatementEnd
