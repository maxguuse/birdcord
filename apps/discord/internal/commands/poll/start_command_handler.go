package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (p *CommandHandler) startPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	var err error
	defer func() {
		if err != nil {
			p.Log.Error("error creating poll", slog.String("error", err.Error()))
			err := interactionRespondError(
				"Произошла ошибка при создании опроса",
				err, p.Session, i,
			)
			if err != nil {
				p.Log.Error(
					"error editing an interaction",
					slog.String("error", err.Error()),
				)
			}

			return
		}

		err = interactionRespondSuccess(
			"Опрос создан!",
			p.Session, i,
		)
		if err != nil {
			p.Log.Error(
				"error editing an interaction",
				slog.String("error", err.Error()),
			)
		}
	}()

	ctx := context.Background()

	err = interactionRespondLoading(
		"Опрос формируется...",
		p.Session, i,
	)
	if err != nil {
		p.Log.Error(
			"error responding to interaction",
			slog.String("error", err.Error()),
		)
		return
	}

	rawOptions := options["options"].StringValue()
	optionsList := strings.Split(rawOptions, "|")
	if len(optionsList) < 2 || len(optionsList) > 25 {
		err = errors.Join(
			domain.ErrUserSide,
			domain.ErrWrongPollOptionsAmount,
		)
		return
	}
	if lo.SomeBy(optionsList, func(o string) bool {
		return utf8.RuneCountInString(o) > 50
	}) {
		err = errors.Join(
			domain.ErrUserSide,
			domain.ErrWrongPollOptionLength,
		)
		return
	}

	guild, err := p.Database.Guilds().GetGuildByDiscordID(ctx, i.GuildID)
	if err != nil {
		return
	}

	user, err := p.Database.Users().GetUserByDiscordID(ctx, i.Member.User.ID)
	if err != nil {
		return
	}

	poll, err := p.Database.Polls().CreatePoll(
		ctx,
		options["title"].StringValue(),
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return
	}

	buttons := make([]discordgo.MessageComponent, 0, len(poll.Options))
	for i, option := range poll.Options {
		customId := fmt.Sprintf("poll_%d_option_%d", poll.ID, option.ID)
		buttons = append(buttons, discordgo.Button{
			Label:    option.Title,
			Style:    discordgo.PrimaryButton,
			CustomID: customId,
		})

		p.EventBus.Subscribe(customId, &VoteButtonHandler{
			poll_id:   int32(poll.ID),
			option_id: int32(option.ID),
			Log:       p.Log,
			Database:  p.Database,
			Session:   p.Session,
		})

		optionsList[i] = fmt.Sprintf("**%d.** %s", i+1, option.Title)
	}
	buttonsGroups := lo.Chunk(buttons, 5)
	actionRows := lo.Map(buttonsGroups, func(buttons []discordgo.MessageComponent, _ int) discordgo.MessageComponent {
		return discordgo.ActionsRow{
			Components: buttons,
		}
	})

	pollEmbed := buildPollEmbed(poll, optionsList, i.Member.User, 0)
	msg, err := p.Session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embeds:     pollEmbed,
		Components: actionRows,
	})
	if err != nil {
		err = errors.Join(domain.ErrInternal, err)
		return
	}

	_, err = p.Database.Polls().CreatePollMessage(
		ctx,
		msg.ID, msg.ChannelID,
		poll.ID,
	)
	if err != nil {
		deleteErr := p.Session.ChannelMessageDelete(i.ChannelID, msg.ID)
		err = errors.Join(domain.ErrInternal, deleteErr, err)
		return
	}
}
