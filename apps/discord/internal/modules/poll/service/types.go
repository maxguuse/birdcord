package service

import "github.com/maxguuse/birdcord/apps/discord/internal/domain"

type GetPollRequest struct {
	GuildID string
	UserID  string
	PollID  int64
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
