package client

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

func (c *Client) onStatusChanged(_ *discordgo.Session, u *discordgo.PresenceUpdate) {
	if u.User.Bot {
		return
	}

	if c.cfg.Environment != config.EnvProduction && u.GuildID != c.cfg.DebugGuildId {
		return
	}

	ctx := context.Background()
	logFieldsFunc := func(err error) []any {
		return []any{
			slog.String("guild_id", u.GuildID),
			slog.String("user_id", u.User.ID),
			slog.Any("error", err),
		}
	}

	go c.handleStreaming(ctx, logFieldsFunc, u)
}

func (c *Client) handleStreaming(
	ctx context.Context,
	logFieldsFunc func(err error) []any,
	u *discordgo.PresenceUpdate,
) {
	isStreaming := lo.SomeBy(u.Activities, func(a *discordgo.Activity) bool {
		return a.Type == discordgo.ActivityTypeStreaming
	})

	_, err := c.redisC.Get(ctx, u.User.ID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		c.logger.Error("could not get streaming user from redis", logFieldsFunc(err)...)

		return
	}

	if !isStreaming && errors.Is(err, redis.Nil) {
		return
	}

	req := &service.SwapUserLiverolesRequest{
		Session: c.router.Session(),
		GuildID: u.GuildID,
		UserID:  u.User.ID,
	}

	if isStreaming {
		if err := c.lrserv.GiveLiveroles(ctx, req); err != nil {
			c.logger.Error("could not give liveroles", logFieldsFunc(err)...)
		}

		if err := c.redisC.Set(ctx, u.User.ID, "streaming", 0).Err(); err != nil {
			c.logger.Error("could not save streaming user in redis", logFieldsFunc(err)...)
		}

		return
	}

	if err := c.lrserv.WithdrawLiveroles(ctx, req); err != nil {
		c.logger.Error("could not withdraw liveroles", logFieldsFunc(err)...)
	}

	if err := c.redisC.Del(ctx, u.User.ID).Err(); err != nil {
		c.logger.Error("could not remove streaming user from redis", logFieldsFunc(err)...)
	}
}
