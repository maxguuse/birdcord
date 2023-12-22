package session

import (
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/libs/config"
)

func New(cfg *config.Config) *discordgo.Session {
	s, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		panic(err)
	}

	return s
}
