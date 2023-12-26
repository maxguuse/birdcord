module github.com/maxguuse/birdcord/libs/migrations

go 1.21.5

require (
	github.com/jackc/pgx/v5 v5.5.0
	github.com/maxguuse/birdcord/libs/config v0.0.0-00010101000000-000000000000
	github.com/pressly/goose/v3 v3.16.0
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/maxguuse/birdcord/libs/config => ../../libs/config
