package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"

	"github.com/maxguuse/birdcord/apps/discord/internal/modules"
	lrservice "github.com/maxguuse/birdcord/apps/discord/internal/modules/liverole/service"
	"github.com/maxguuse/birdcord/libs/config"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/maxguuse/disroute"
	"go.uber.org/fx"
)

type Client struct {
	router *disroute.Router
	logger logger.Logger
	cfg    *config.Config

	redisC *redis.Client
	lrserv *lrservice.Service
}

type ClientOpts struct {
	fx.In
	LC             fx.Lifecycle
	CommandHandler *modules.Handler

	Cfg    *config.Config
	Logger logger.Logger

	Redis         *redis.Client
	LiverolesServ *lrservice.Service
}

func New(opts ClientOpts) error {
	router, err := disroute.New(opts.Cfg.DiscordToken)
	if err != nil {
		return err
	}

	router.Session().Identify.Intents = discordgo.IntentsAll

	c := &Client{
		router: router,
		logger: opts.Logger,
		cfg:    opts.Cfg,

		redisC: opts.Redis,
		lrserv: opts.LiverolesServ,
	}

	if opts.Cfg.Environment != config.EnvProduction {
		router.Use(func(hf disroute.HandlerFunc) disroute.HandlerFunc {
			return func(ctx *disroute.Ctx) disroute.Response {
				c.logger.Debug("interaction",
					slog.String("type", ctx.Interaction().Type.String()),
					slog.String("id", ctx.Interaction().ID),
					slog.String("user", ctx.Interaction().Member.User.GlobalName))

				return hf(ctx)
			}
		})
	}

	router.Use(func(hf disroute.HandlerFunc) disroute.HandlerFunc {
		return func(ctx *disroute.Ctx) disroute.Response {
			err := ctx.Session().InteractionRespond(ctx.Interaction(), &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})

			if err != nil {
				c.logger.Error("error responding to interaction", slog.Any("error", err))
			}

			return hf(ctx)
		}
	})

	router.SetResponseHandler(c.responseHandler)
	router.SetComponentKeyFunc(c.componentKeyFunc)

	c.registerLogger()
	c.registerHandlers()

	opts.CommandHandler.Register(router)

	opts.LC.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return c.router.Open()
		},
		OnStop: func(_ context.Context) error {
			return c.router.Session().Close()
		},
	})

	return nil
}

var NewFx = fx.Options(
	fx.Invoke(
		New,
	),
)
