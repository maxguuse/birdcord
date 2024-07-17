package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	. "github.com/maxguuse/birdcord/libs/jet/generated/birdcord/public/table"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TxGetter *txmanager.TxGetter
}

func NewPgx(opts Opts) *pollsPgx {
	return &pollsPgx{
		txGetter: opts.TxGetter,
	}
}

var _ Repository = &pollsPgx{}

type pollsPgx struct {
	txGetter *txmanager.TxGetter
}

func (p *pollsPgx) CreatePoll(
	ctx context.Context,
	discordGuildId, discordAuthorId int64,
	title string,
	pollOptions []string,
) (*domain.PollWithDetails, error) {
	dest := &domain.PollWithDetails{}

	insertedPoll := postgres.CTE("inserted_poll")
	insertedPollOptions := postgres.CTE("inserted_poll_options")
	err := postgres.WITH(
		insertedPoll.AS(
			Polls.INSERT(
				Polls.Title,
				Polls.AuthorID,
				Polls.GuildID,
			).VALUES(
				title,
				discordAuthorId,
				discordGuildId,
			).RETURNING(Polls.AllColumns),
		),
		insertedPollOptions.AS(
			PollOptions.INSERT(
				PollOptions.Title,
				PollOptions.PollID,
			).VALUES(
				postgres.Raw("UNNEST($1::varchar[])", map[string]any{
					"$1": encodeStringsSlice(pollOptions),
				}),
				insertedPoll.SELECT(Polls.ID.From(insertedPoll)),
			).RETURNING(PollOptions.AllColumns),
		),
	)(
		postgres.SELECT(
			insertedPoll.AllColumns(),
			insertedPollOptions.AllColumns(),
		).FROM(
			insertedPoll.LEFT_JOIN(
				insertedPollOptions,
				Polls.ID.From(insertedPoll).EQ(PollOptions.PollID.From(insertedPollOptions)),
			),
		),
	).QueryContext(ctx, p.txGetter.DefaultTxOrDB(ctx), dest)
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

	err := postgres.SELECT(
		Polls.AllColumns,
		PollOptions.AllColumns,
		PollMessages.AllColumns,
		PollVotes.AllColumns,
	).FROM(
		Polls.LEFT_JOIN(
			PollOptions,
			Polls.ID.EQ(PollOptions.PollID),
		).LEFT_JOIN(
			PollMessages,
			Polls.ID.EQ(PollMessages.PollID),
		).LEFT_JOIN(
			PollVotes,
			Polls.ID.EQ(PollVotes.PollID),
		),
	).WHERE(
		Polls.ID.EQ(postgres.Int(int64(pollId))),
	).QueryContext(ctx, p.txGetter.DefaultTxOrDB(ctx), dest)
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return dest, nil
}

func (p *pollsPgx) GetActivePolls(
	ctx context.Context,
	discordGuildId int64, discordAuthorId int64,
) ([]*domain.Poll, error) {
	dest := []*domain.Poll{}

	err := postgres.SELECT(
		Polls.ID,
		Polls.Title,
		Polls.IsActive,
		Polls.CreatedAt,
	).FROM(
		Polls,
	).WHERE(
		Polls.IsActive.EQ(postgres.Bool(true)).
			AND(Polls.GuildID.EQ(postgres.Int64(discordGuildId))).
			AND(Polls.AuthorID.EQ(postgres.Int64(discordAuthorId))),
	).QueryContext(ctx, p.txGetter.DefaultTxOrDB(ctx), &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (p *pollsPgx) TryAddVote(
	ctx context.Context,
	discordUserId int64,
	pollId, optionId int,
) (*domain.PollVote, error) {
	dest := &domain.PollVote{}

	err := PollVotes.INSERT(
		PollVotes.UserID,
		PollVotes.PollID,
		PollVotes.OptionID,
	).VALUES(
		discordUserId,
		pollId,
		optionId,
	).RETURNING(
		PollVotes.AllColumns,
	).QueryContext(ctx, p.txGetter.DefaultTxOrDB(ctx), dest)
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return dest, nil
}

func (p *pollsPgx) CreatePollMessage(
	ctx context.Context,
	discordMessageId, discordChannelId int64,
	pollId int,
) (*domain.PollMessage, error) {
	dest := &domain.PollMessage{}

	err := PollMessages.INSERT(
		PollMessages.PollID,
		PollMessages.DiscordMessageID,
		PollMessages.DiscordChannelID,
	).VALUES(
		pollId,
		discordMessageId,
		discordChannelId,
	).RETURNING(
		PollMessages.AllColumns,
	).QueryContext(ctx, p.txGetter.DefaultTxOrDB(ctx), dest)
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return dest, nil
}

func (p *pollsPgx) UpdatePollStatus(
	ctx context.Context,
	pollId int,
	isActive bool,
) error {
	_, err := Polls.UPDATE(
		Polls.IsActive,
	).SET(
		isActive,
	).WHERE(
		Polls.ID.EQ(postgres.Int(int64(pollId))),
	).ExecContext(ctx, p.txGetter.DefaultTxOrDB(ctx))
	if err != nil {
		return err // TODO: Wrap error
	}

	return nil
}

func (p *pollsPgx) AddPollOption(
	ctx context.Context,
	pollId int,
	pollOption string,
) (*domain.PollOption, error) {
	dest := &domain.PollOption{}

	err := PollOptions.INSERT(
		PollOptions.Title,
		PollOptions.PollID,
	).VALUES(
		pollOption,
		pollId,
	).RETURNING(
		PollOptions.AllColumns,
	).QueryContext(ctx, p.txGetter.DefaultTxOrDB(ctx), dest)
	if err != nil {
		return nil, err // TODO: Wrap error
	}

	return dest, nil
}

func (p *pollsPgx) RemovePollOption(
	ctx context.Context,
	optionId int,
) error {
	_, err := PollOptions.
		DELETE().
		WHERE(
			PollOptions.ID.EQ(postgres.Int(int64(optionId))),
		).ExecContext(ctx, p.txGetter.DefaultTxOrDB(ctx))
	if err != nil {
		return err // TODO: Wrap error
	}

	return nil
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
