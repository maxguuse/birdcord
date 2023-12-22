package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
	"github.com/samber/lo"
)

func (p *CommandHandler) startPoll(
	i *discordgo.Interaction,
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
) {
	ctx := context.Background()

	err := interactionRespondLoading(
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

	err = p.Database.Transaction(func(q *queries.Queries) error {
		guild, getGuildErr := q.GetGuildByDiscordID(ctx, i.GuildID)
		if getGuildErr != nil {
			return getGuildErr
		}

		user, getUserErr := q.GetUserByDiscordID(ctx, i.Member.User.ID)
		if getUserErr != nil && !errors.Is(getUserErr, pgx.ErrNoRows) {
			return getUserErr
		}
		if user.ID == 0 {
			var createUserErr error
			user, createUserErr = q.CreateUser(ctx, i.Member.User.ID)
			if createUserErr != nil {
				return createUserErr
			}
		}

		poll, createPolLErr := q.CreatePoll(ctx, queries.CreatePollParams{
			Title: options["title"].StringValue(),
			AuthorID: pgtype.Int4{
				Int32: user.ID,
				Valid: true,
			},
			GuildID: guild.ID,
		})
		if createPolLErr != nil {
			return createPolLErr
		}

		rawOptions := options["options"].StringValue()
		optionsList := strings.Split(rawOptions, "|")
		if len(optionsList) < 2 {
			return fmt.Errorf("not enough options")
		}

		buttons := make([]discordgo.MessageComponent, 0, len(optionsList))
		for i, option := range optionsList {
			if len(option) == 0 || utf8.RuneCountInString(option) > 50 {
				return fmt.Errorf("invalid option length")
			}

			pollOption, err := q.CreatePollOption(ctx, queries.CreatePollOptionParams{
				Title:  option,
				PollID: poll.ID,
			})
			if err != nil {
				return err //TODO add ErrQueryFailed in db lib
			}

			customId := fmt.Sprintf("poll_%d_option_%d", poll.ID, pollOption.ID)
			buttons = append(buttons, discordgo.Button{
				Label:    pollOption.Title,
				Style:    discordgo.PrimaryButton,
				CustomID: customId,
			})

			p.EventBus.Subscribe(customId, &VoteButtonHandler{
				poll_id:   poll.ID,
				option_id: pollOption.ID,
				Log:       p.Log,
				Database:  p.Database,
				Session:   p.Session,
			})

			optionsList[i] = fmt.Sprintf("**%d.** %s", i+1, option)
		}
		buttonsGroups := lo.Chunk(buttons, 5)
		actionRows := lo.Map(buttonsGroups, func(buttons []discordgo.MessageComponent, _ int) discordgo.MessageComponent {
			return discordgo.ActionsRow{
				Components: buttons,
			}
		})

		a := buildPollEmbed(
			poll,
			optionsList,
			i.Member.User,
			0,
		)
		discordMsg, sendMessageErr := p.Session.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Embeds:     a,
			Components: actionRows,
		})
		if sendMessageErr != nil {
			return sendMessageErr
		}

		msg, createMessageErr := q.CreateMessage(ctx, queries.CreateMessageParams{
			DiscordMessageID: discordMsg.ID,
			DiscordChannelID: discordMsg.ChannelID,
		})
		if createMessageErr != nil {
			deleteMessageErr := p.Session.ChannelMessageDelete(i.ChannelID, discordMsg.ID)
			return errors.Join(createMessageErr, deleteMessageErr)
		}

		createPollMessageErr := q.CreatePollMessage(ctx, queries.CreatePollMessageParams{
			PollID:    poll.ID,
			MessageID: msg.ID,
		})
		if createPollMessageErr != nil {
			deleteMessageErr := p.Session.ChannelMessageDelete(i.ChannelID, discordMsg.ID)
			return errors.Join(createPollMessageErr, deleteMessageErr)
		}

		return nil
	})

	var pgErr *pgconn.PgError

	if err != nil && errors.Is(err, pgErr) {
		p.Log.Error("error creating poll", slog.String("error", err.Error()))
		err := interactionRespondError(
			"Произошла внутренняя ошибка при формировании опроса",
			fmt.Errorf("internal error"),
			p.Session, i,
		)
		if err != nil {
			p.Log.Error(
				"error editing an interaction",
				slog.String("error", err.Error()),
			)
		}

		return
	}

	if err != nil {
		p.Log.Error("error creating poll", slog.String("error", err.Error()))
		err := interactionRespondError(
			"Произошла ошибка при формировании опроса",
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
}
