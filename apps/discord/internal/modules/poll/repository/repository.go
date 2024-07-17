package repository

import (
	"context"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
)

type Repository interface {
	CreatePoll(
		ctx context.Context,
		discordGuildId, discordAuthorId int64,
		title string, pollOptions []string,
	) (*domain.PollWithDetails, error)
	GetPollWithDetails(
		ctx context.Context,
		pollId int,
	) (*domain.PollWithDetails, error)
	TryAddVote(
		ctx context.Context,
		discordUserId int64,
		pollId, optionId int,
	) (*domain.PollVote, error)
	CreatePollMessage(
		ctx context.Context,
		discordMessageId, discordChannelId int64,
		pollId int,
	) (*domain.PollMessage, error)
	UpdatePollStatus(
		ctx context.Context,
		pollId int,
		isActive bool,
	) error
	GetActivePolls(
		ctx context.Context,
		discordGuildId, discordAuthorId int64,
	) ([]*domain.Poll, error)
	AddPollOption(
		ctx context.Context,
		pollId int,
		pollOption string,
	) (*domain.PollOption, error)
	RemovePollOption(
		ctx context.Context,
		optionId int,
	) error
}
