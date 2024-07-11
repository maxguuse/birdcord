package domain

import "time"

type Poll struct {
	ID        int       `alias:"polls.id"`
	Title     string    `alias:"polls.title"`
	IsActive  bool      `alias:"polls.is_active"`
	CreatedAt time.Time `alias:"polls.created_at"`
}

type PollDetails struct {
	Guild  Guild
	Author User

	Options  []PollOption
	Messages []PollMessage
	Votes    []PollVote
}

type PollWithDetails struct {
	Poll
	PollDetails
}

type PollOption struct {
	ID    int    `alias:"poll_options.id"`
	Title string `alias:"poll_options.title"`
}

type PollMessage struct {
	ID               int    `alias:"poll_messages.id"`
	MessageID        int    `alias:"poll_messages.message_id"`
	DiscordMessageID string `alias:"messages.discord_message_id"`
	DiscordChannelID string `alias:"messages.discord_channel_id"`
}

type PollVote struct {
	ID       int `alias:"poll_votes.id"`
	OptionID int `alias:"poll_votes.option_id"`
	UserID   int `alias:"poll_votes.user_id"`
}

type User struct {
	ID            int    `alias:"users.id"`
	DiscordUserID string `alias:"users.discord_user_id"`
}

type Guild struct {
	ID             int    `alias:"guilds.id"`
	DiscordGuildID string `alias:"guilds.discord_guild_id"`
}

type Liverole struct {
	ID            int
	GuildID       int
	RoleID        int
	DiscordRoleID string
}
