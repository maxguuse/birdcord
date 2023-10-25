module github.com/maxguuse/birdcord/apps/bot

go 1.21.3

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/maxguuse/birdcord/libs/types v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.20.1
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/maxguuse/birdcord/libs/types => ../../libs/types
