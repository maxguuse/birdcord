package service

import "github.com/maxguuse/birdcord/apps/discord/internal/domain"

type GetPollRequest struct {
	GuildID string
	UserID  string
	PollID  int64
}

type GetActivePollsRequest struct {
	GuildID string
	UserID  string
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

type StopRequest struct {
	GuildID string
	UserID  string
	PollID  int64
}

type StopResponse struct {
	Poll    *domain.PollWithDetails
	Winners []string
}

type AddOptionRequest struct {
	GuildID string
	UserID  string
	PollID  int64
	Option  string
}

type RemoveOptionRequest struct {
	GuildID  string
	UserID   string
	PollID   int64
	OptionID int64
}

type AddVoteRequest struct {
	UserID   string
	CustomID string
}

type CreateMessageRequest struct {
	PollID  int
	Message Message
}

type Message struct {
	ID        string
	ChannelID string
}
