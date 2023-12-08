module github.com/maxguuse/birdcord/apps/discord

go 1.21.5

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/bwmarrin/discordgo v0.27.1
	github.com/jackc/pgx/v5 v5.5.0
	github.com/maxguuse/birdcord/libs/config v0.0.0-00010101000000-000000000000
	github.com/maxguuse/birdcord/libs/logger v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.39.0
	go.uber.org/fx v1.20.1
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/exp v0.0.0-20231206192017-f3f8817b8deb // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/maxguuse/birdcord/libs/logger => ../../libs/logger

replace github.com/maxguuse/birdcord/libs/config => ../../libs/config
