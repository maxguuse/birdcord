package domain

import "time"

type Poll struct {
	ID        int       `sql:"primary_key" alias:"polls.id"`
	Title     string    `alias:"polls.title"`
	IsActive  bool      `alias:"polls.is_active"`
	CreatedAt time.Time `alias:"polls.created_at"`
}

type PollDetails struct {
	GuildID  int `alias:"polls.guild_id"`
	AuthorID int `alias:"polls.author_id"`

	Options  []PollOption
	Messages []PollMessage
	Votes    []PollVote
}

type PollWithDetails struct {
	Poll
	PollDetails
}

type PollOption struct {
	ID    int    `sql:"primary_key" alias:"poll_options.id"`
	Title string `alias:"poll_options.title"`
}

type PollMessage struct {
	ID               int    `sql:"primary_key" alias:"poll_messages.id"`
	MessageID        int    `alias:"poll_messages.message_id"`
	DiscordMessageID string `alias:"poll_messages.discord_message_id"`
	DiscordChannelID string `alias:"poll_messages.discord_channel_id"`
}

type PollVote struct {
	ID       int `sql:"primary_key" alias:"poll_votes.id"`
	OptionID int `alias:"poll_votes.option_id"`
	UserID   int `alias:"poll_votes.user_id"`
}

type Liverole struct {
	ID             int `sql:"primary_key" alias:"liveroles.id"`
	DiscordRoleID  int `alias:"liveroles.discord_role_id"`
	DiscordGuildID int `alias:"liveroles.discord_guild_id"`
}
