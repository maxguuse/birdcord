package domain

import "time"

type Poll struct {
	ID        int
	Title     string
	CreatedAt time.Time
}

type PollDetails struct {
	Guild  PollGuild
	Author PollAuthor

	Options  []PollOption
	Messages []PollMessage
	Votes    []PollVote
}

type PollWithDetails struct {
	Poll
	PollDetails
}

type PollGuild struct {
	ID             int
	DiscordGuildID string
}

type PollAuthor struct {
	ID            int
	DiscordUserID string
}

type PollOption struct {
	ID    int
	Title string
}

type PollMessage struct {
	ID               int
	MessageID        int
	DiscordMessageID string
	DiscordChannelID string
}

type PollVote struct {
	ID       int
	OptionID int
	UserID   int
}

type User struct {
	ID            int
	DiscordUserID string
}

type Guild struct {
	ID             int
	DiscordGuildID string
}
