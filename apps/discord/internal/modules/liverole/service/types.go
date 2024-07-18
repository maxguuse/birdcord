package service

import "github.com/bwmarrin/discordgo"

type AddLiveRoleRequest struct {
	GuildID string
	RoleID  string
}

type RemoveLiveRoleRequest struct {
	GuildID string
	RoleID  string
}

type SwapUserLiverolesRequest struct {
	Session *discordgo.Session
	GuildID string
	UserID  string
}
