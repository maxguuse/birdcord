package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/bot/internal/scommands/polls"
	"go.uber.org/fx"
	"os"
)

type Bot struct {
	session *discordgo.Session
	polls   *polls.Polls
}

func New(lc fx.Lifecycle, p *polls.Polls) (*Bot, error) {
	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		session: s,
		polls:   p,
	}

	bot.session.Identify.Presence.Game = discordgo.Activity{
		Name: "как Гусь кодит",
		Type: discordgo.ActivityTypeWatching,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := s.Open()
			if err != nil {
				return fmt.Errorf("error opening connection: %v", err)
			}

			bot.SetupHandlers()
			bot.SetupIntents()
			bot.SetupScommands()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := s.Close()
			if err != nil {
				return fmt.Errorf("error closing connection: %v", err)
			}
			return nil
		},
	})

	return bot, nil
}
