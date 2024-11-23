-- +goose Up
-- +goose StatementBegin
select 'up SQL query'
;

CREATE TABLE temp_voice_hub (
	id serial,
	discord_channel_id bigint,
	discord_guild_id bigint,
	tempvoice_template varchar,
	tempvoice_category bigint
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
select 'down SQL query'
;

DROP TABLE temp_voice_hub;
-- +goose StatementEnd


