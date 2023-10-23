FROM golang:1.21.3-alpine as builder
WORKDIR /app
RUN apk add upx
COPY . .

FROM builder as bot_builder
RUN cd apps/bot && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM alpine:latest as bot
WORKDIR /app
COPY --from=bot_builder /app/apps/bot/out /bin/bot
CMD ["/bin/bot"]
