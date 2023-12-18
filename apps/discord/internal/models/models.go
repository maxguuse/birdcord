package models

import "time"

type Guild struct {
	ID        string
	DiscordID string
}

type User struct {
	ID        string
	DiscordID string
}

type Poll struct {
	ID        string
	Title     string
	CreatedAt time.Time
	GuildID   string
	AuthorID  string
}
