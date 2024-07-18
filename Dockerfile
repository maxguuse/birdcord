FROM golang:1.22.2-alpine as builder
WORKDIR /app
RUN apk add upx
COPY . .
RUN go work sync

FROM builder as discord_builder
RUN cd apps/discord && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM alpine:latest as discord
WORKDIR /app
COPY --from=discord_builder /app/apps/discord/out /bin/discord
CMD ["/bin/discord"]

FROM builder as migrations_builder
RUN cd libs/migrations && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./main.go && upx -9 -k ./out

FROM alpine:latest as migrations
WORKDIR /app
COPY --from=migrations_builder /app/libs/migrations/out /bin/migrations
CMD ["/bin/migrations"]