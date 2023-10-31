module github.com/maxguuse/birdcord/apps/bot

go 1.21.3

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/maxguuse/birdcord/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.38.1
	go.uber.org/fx v1.20.1
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
)

replace github.com/maxguuse/birdcord/libs/grpc => ../../libs/grpc
