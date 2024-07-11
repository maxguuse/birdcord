package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/libs/jet/generated/birdcord/public/model"
	. "github.com/maxguuse/birdcord/libs/jet/generated/birdcord/public/table"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TxManager *txmanager.TxManager
}

func NewPgx(opts Opts) *pollsPgx {
	return &pollsPgx{
		txm: opts.TxManager,
	}
}

var _ Repository = &pollsPgx{}

type pollsPgx struct {
	txm *txmanager.TxManager
}

type CreatePollDest struct {
	Poll model.Polls

	Guild  model.Guilds
	Author model.Users

	Options []model.PollOptions
}

func (p *pollsPgx) CreatePoll(
	ctx context.Context,
	discordGuildId, discordAuthorId string,
	title string,
	pollOptions []string,
) (*domain.PollWithDetails, error) {
	dest := &domain.PollWithDetails{}

	err := p.txm.Do(ctx, func(db qrm.DB) error {
		insertedPoll := postgres.CTE("inserted_poll")
		err := postgres.WITH(
			insertedPoll.AS(
				Polls.INSERT(
					Polls.Title,
					Polls.AuthorID,
					Polls.GuildID,
				).VALUES(
					title,
					postgres.
						SELECT(Users.ID).FROM(Users).
						WHERE(Users.DiscordUserID.EQ(postgres.String(discordAuthorId))),
					postgres.
						SELECT(Guilds.ID).FROM(Guilds).
						WHERE(Guilds.DiscordGuildID.EQ(postgres.String(discordGuildId))),
				).RETURNING(Polls.AllColumns),
			),
		)(
			postgres.SELECT(
				insertedPoll.AllColumns(),
				Guilds.AllColumns,
				Users.AllColumns,
			).FROM(
				insertedPoll.LEFT_JOIN(
					Guilds,
					Polls.GuildID.From(insertedPoll).EQ(Guilds.ID),
				).LEFT_JOIN(
					Users,
					Polls.AuthorID.From(insertedPoll).EQ(Users.ID),
				),
			),
		).QueryContext(ctx, db, dest)
		if err != nil {
			return err // TODO: Wrap error
		}

		err = PollOptions.INSERT(
			PollOptions.Title,
			PollOptions.PollID,
		).VALUES(
			postgres.Raw("UNNEST($1::varchar[])", map[string]any{
				"$1": encodeStringsSlice(pollOptions),
			}),
			dest.Poll.ID,
		).RETURNING(
			PollOptions.AllColumns.Except(PollOptions.PollID),
		).QueryContext(ctx, db, &dest.Options)
		if err != nil {
			return err // TODO: Wrap error
		}

		return nil
	})
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return dest, nil
}

func (p *pollsPgx) GetPollWithDetails(
	ctx context.Context,
	pollId int,
) (*domain.PollWithDetails, error) {
	dest := &domain.PollWithDetails{}

	err := p.txm.Do(ctx, func(db qrm.DB) error {
		err := postgres.SELECT(
			Polls.AllColumns,
			Guilds.AllColumns,
			Users.AllColumns,
			PollOptions.AllColumns,
			PollMessages.AllColumns,
			Messages.AllColumns,
			PollVotes.AllColumns,
		).FROM(
			Polls.LEFT_JOIN(
				Guilds,
				Polls.GuildID.EQ(Guilds.ID),
			).LEFT_JOIN(
				Users,
				Polls.AuthorID.EQ(Users.ID),
			).LEFT_JOIN(
				PollOptions,
				Polls.ID.EQ(PollOptions.PollID),
			).LEFT_JOIN(
				PollMessages,
				Polls.ID.EQ(PollMessages.PollID),
			).LEFT_JOIN(
				Messages,
				PollMessages.MessageID.EQ(Messages.ID),
			).LEFT_JOIN(
				PollVotes,
				Polls.ID.EQ(PollVotes.PollID),
			),
		).WHERE(
			Polls.ID.EQ(postgres.Int(int64(pollId))),
		).QueryContext(ctx, db, dest)
		if err != nil {
			return err // TODO: Wrap error
		}

		return nil
	})
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return dest, nil
}

func (p *pollsPgx) TryAddVote(
	ctx context.Context,
	discordUserId string,
	pollId, optionId int,
) (*domain.PollVote, error) {
	panic("TODO: Implement")
}

func (p *pollsPgx) CreatePollMessage(
	ctx context.Context,
	discordMessageId, discordChannelId string,
	pollId int,
) (*domain.PollMessage, error) {
	panic("TODO: Implement")
}

func (p *pollsPgx) UpdatePollStatus(
	ctx context.Context,
	pollId int,
	isActive bool,
) error {
	panic("TODO: Implement")
}

func (p *pollsPgx) GetActivePolls(
	ctx context.Context,
	discordGuildId, discordAuthorId string,
) ([]*domain.Poll, error) {
	panic("TODO: Implement")
}

func (p *pollsPgx) AddPollOption(
	ctx context.Context,
	pollId int,
	pollOption string,
) (*domain.PollOption, error) {
	panic("TODO: Implement")
}

func (p *pollsPgx) RemovePollOption(
	ctx context.Context,
	optionId int,
) error {
	panic("TODO: Implement")
}

var quoteArrayReplacer = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func quoteArrayElement(src string) string {
	return `"` + quoteArrayReplacer.Replace(src) + `"`
}

func encodeStringsSlice(strs []string) string {
	for i, str := range strs {
		strs[i] = quoteArrayElement(str)
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ","))
}
