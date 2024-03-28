package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/repository"
	"github.com/samber/lo"
)

type Service struct {
	db repository.DB
}

func New(db repository.DB) *Service {
	return &Service{
		db: db,
	}
}

type CreateRequest struct {
	GuildID string
	UserID  string
	Poll    Poll
}

type Poll struct {
	Title   string
	Options string
}

func (s *Service) Create(ctx context.Context, r *CreateRequest) (*domain.PollWithDetails, error) {
	optionsList, err := processPollOptions(r.Poll.Options)
	if err != nil {
		return nil, err
	}

	guild, err := s.db.Guilds().GetGuildByDiscordID(ctx, r.GuildID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	user, err := s.db.Users().GetUserByDiscordID(ctx, r.UserID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	poll, err := s.db.Polls().CreatePoll(
		ctx,
		r.Poll.Title,
		guild.ID,
		user.ID,
		optionsList,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, err)
	}

	return poll, nil
}

func processPollOptions(rawOptions string) ([]string, error) {
	optionsList := strings.Split(rawOptions, "|")
	if len(optionsList) < 2 || len(optionsList) > 25 {
		return nil, &domain.UsersideError{
			Msg: "Количество вариантов опроса должно быть от 2 до 25 включительно.",
		}
	}
	if lo.SomeBy(optionsList, func(o string) bool {
		return utf8.RuneCountInString(o) > 50 || utf8.RuneCountInString(o) < 1
	}) {
		return nil, &domain.UsersideError{
			Msg: "Длина варианта опроса не может быть больше 50 или меньше 1 символа.",
		}
	}

	return optionsList, nil
}
