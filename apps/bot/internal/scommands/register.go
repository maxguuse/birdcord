package scommands

import (
	"github.com/bwmarrin/discordgo"
)

func Register(s *discordgo.Session) {
	registerPollsCommands(s)
}
