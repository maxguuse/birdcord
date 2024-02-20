package client

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (c *Client) registerHandlers() {
	c.AddHandler(c.onInteractionCreate)
	c.AddHandler(c.onMessageDelete)
	c.AddHandler(c.onReady)
	c.AddHandler(c.onConnect)
	c.AddHandler(c.onDisconnect)
	c.AddHandler(c.onStatusChanged)
}

func (c *Client) onConnect(_ *discordgo.Session, _ *discordgo.Connect) {
	c.Log.Info("Bot is connected!")
}

func (c *Client) onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	c.Log.Info("Bot is disconnected!")
}

func (c *Client) onStatusChanged(_ *discordgo.Session, u *discordgo.PresenceUpdate) {
	if u.User.Bot {
		return
	}

	ctx := context.Background()

	isStreaming := lo.SomeBy(u.Activities, func(a *discordgo.Activity) bool {
		return a.Type == discordgo.ActivityTypeStreaming
	})

	guild, err := c.Database.Guilds().GetGuildByDiscordID(ctx, u.GuildID)
	if err != nil {
		c.Log.Error("could not get guild", err)
	}

	roles, err := c.Database.Liveroles().GetLiveroles(ctx, guild.ID)
	if err != nil {
		c.Log.Error("could not give streaming role", err)
	}

	member, err := c.GuildMember(u.GuildID, u.User.ID)
	if err != nil {
		c.Log.Error("could not get member", err)
	}

	liverolesIds := lo.Map(roles, func(role *domain.Liverole, _ int) string { return role.DiscordRoleID })
	memberRolesIds := member.Roles

	if isStreaming {
		for _, role := range roles {
			err = c.GuildMemberRoleAdd(u.GuildID, u.User.ID, role.DiscordRoleID)
			if err != nil {
				c.Log.Error("could not add role", err)
			}
			c.Log.Debug("Role added",
				slog.String("role", role.DiscordRoleID),
				slog.String("user", u.Presence.User.ID),
			)
		}
	} else {
		for _, liveroleId := range liverolesIds {
			if !isStreaming && lo.Contains(memberRolesIds, liveroleId) {
				err = c.GuildMemberRoleRemove(u.GuildID, u.User.ID, liveroleId)
				if err != nil {
					c.Log.Error("could not remove role", err)
				}
			}
		}
	}

	c.Log.Debug("Status changed",
		slog.String("user", u.Presence.User.ID),
		slog.String("status", string(u.Status)),
	)
}
