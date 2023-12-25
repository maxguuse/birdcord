package poll

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/eventbus"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/maxguuse/birdcord/libs/logger"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type AutocompleteHandler struct {
	Log      logger.Logger
	Database repository.DB
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

type AutocompleteHandlerOpts struct {
	fx.In

	Log      logger.Logger
	Database repository.DB
	EventBus *eventbus.EventBus
	Session  *discordgo.Session
}

func NewAutocompleteHandler(opts AutocompleteHandlerOpts) *AutocompleteHandler {
	return &AutocompleteHandler{
		Log:      opts.Log,
		Database: opts.Database,
		EventBus: opts.EventBus,
		Session:  opts.Session,
	}
}

func (p *AutocompleteHandler) Handle(i any) {
	cmd, ok := i.(*discordgo.Interaction)
	if !ok {
		return
	}

	commandOptions := buildCommandOptionsMap(cmd)

	switch cmd.ApplicationCommandData().Options[0].Name {
	case "stop":
		p.autocompletePollList(cmd, commandOptions)
	case "status":
		p.autocompletePollList(cmd, commandOptions)
	}
}

func (a *AutocompleteHandler) autocompletePollList(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	ctx := context.Background()

	guild, err := a.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return
	}

	user, err := a.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return
	}

	polls, err := a.Database.Polls().GetActivePolls(ctx, guild.ID, user.ID)
	if err != nil {
		return
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(polls))
	for i, poll := range polls {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("Poll ID: %d | %s", poll.ID, poll.Title),
			Value: poll.ID,
		}
	}

	err = a.Session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: lo.Filter(choices, func(c *discordgo.ApplicationCommandOptionChoice, _ int) bool {
				s, ok := options["poll"].Value.(string)
				if !ok {
					return false
				}

				return strings.Contains(c.Name, s)
			}),
		},
	})
	if err != nil {
		return
	}
}
