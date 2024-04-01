// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package queries

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Guild struct {
	ID             int32  `json:"id"`
	DiscordGuildID string `json:"discord_guild_id"`
}

type Liverole struct {
	ID     int32 `json:"id"`
	RoleID int32 `json:"role_id"`
}

type Message struct {
	ID               int32  `json:"id"`
	DiscordMessageID string `json:"discord_message_id"`
	DiscordChannelID string `json:"discord_channel_id"`
}

type Poll struct {
	ID        int32            `json:"id"`
	Title     string           `json:"title"`
	IsActive  bool             `json:"is_active"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	GuildID   int32            `json:"guild_id"`
	AuthorID  pgtype.Int4      `json:"author_id"`
}

type PollMessage struct {
	ID        int32 `json:"id"`
	MessageID int32 `json:"message_id"`
	PollID    int32 `json:"poll_id"`
}

type PollOption struct {
	ID     int32  `json:"id"`
	Title  string `json:"title"`
	PollID int32  `json:"poll_id"`
}

type PollVote struct {
	ID       int32 `json:"id"`
	PollID   int32 `json:"poll_id"`
	OptionID int32 `json:"option_id"`
	UserID   int32 `json:"user_id"`
}

type Role struct {
	ID            int32  `json:"id"`
	GuildID       int32  `json:"guild_id"`
	DiscordRoleID string `json:"discord_role_id"`
}

type User struct {
	ID            int32  `json:"id"`
	DiscordUserID string `json:"discord_user_id"`
}
