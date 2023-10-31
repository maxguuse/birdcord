FROM golang:1.21.3-alpine as builder
WORKDIR /app
RUN apk add upx
COPY apps apps
COPY libs libs

FROM builder as bot_builder
RUN cd apps/bot && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM alpine:latest as bot
WORKDIR /app
COPY --from=bot_builder /app/apps/bot/out /bin/bot
CMD ["/bin/bot"]

FROM builder as polls_builder
RUN cd apps/polls && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM alpine:latest as polls
WORKDIR /app
COPY --from=polls_builder /app/apps/polls/out /bin/polls
CMD ["/bin/polls"]

FROM builder as migrations_builder
RUN cd libs/migrations && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./main.go && upx -9 -k ./out

FROM alpine:latest as migrations
WORKDIR /app
COPY --from=migrations_builder /app/libs/migrations/out /bin/migrations
CMD ["/bin/migrations"]