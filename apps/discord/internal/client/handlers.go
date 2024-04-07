package client

import (
	"context"
	"log/slog"

	"github.com/avast/retry-go/v4"
	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func (c *Client) registerHandlers() {
	c.router.Session().AddHandler(c.router.InteractionHandler)

	c.router.Session().AddHandler(c.onMessageDelete)
	c.router.Session().AddHandler(c.onReady)
	c.router.Session().AddHandler(c.onConnect)
	c.router.Session().AddHandler(c.onDisconnect)
	c.router.Session().AddHandler(c.onStatusChanged)
}

func (c *Client) onConnect(_ *discordgo.Session, _ *discordgo.Connect) {
	c.logger.Info("Bot is connected!")
}

func (c *Client) onDisconnect(_ *discordgo.Session, _ *discordgo.Disconnect) {
	c.logger.Info("Bot is disconnected!")
}

func (c *Client) onStatusChanged(_ *discordgo.Session, u *discordgo.PresenceUpdate) {
	if u.User.Bot {
		return
	}

	ctx := context.Background()

	isStreaming := lo.SomeBy(u.Activities, func(a *discordgo.Activity) bool {
		return a.Type == discordgo.ActivityTypeStreaming
	})

	guild, err := c.db.Guilds().GetGuildByDiscordID(ctx, u.GuildID)
	if err != nil {
		c.logger.Error("could not get guild", err)
	}

	roles, err := c.db.Liveroles().GetLiveroles(ctx, guild.ID)
	if err != nil {
		c.logger.Error("could not give streaming role", err)
	}

	var member *discordgo.Member
	err = retry.Do(func() error {
		m, err := c.router.Session().GuildMember(u.GuildID, u.User.ID)
		if err != nil {
			return err
		}

		member = m

		return nil
	},
		retry.Attempts(5),
	)
	if err != nil {
		c.logger.Error("failed to fetch member from discord",
			slog.Any("error", err),
			slog.String("guild", u.GuildID),
		)

		return
	}

	liverolesIds := lo.Map(roles, func(role *domain.Liverole, _ int) string { return role.DiscordRoleID })
	memberRolesIds := member.Roles

	if isStreaming {
		for _, role := range roles {
			err = c.router.Session().GuildMemberRoleAdd(u.GuildID, u.User.ID, role.DiscordRoleID)
			if err != nil {
				c.logger.Error("could not add role", err)
			}
			c.logger.Debug("Role added",
				slog.String("role", role.DiscordRoleID),
				slog.String("user", u.Presence.User.ID),
			)
		}
	} else {
		for _, liveroleId := range liverolesIds {
			if !isStreaming && lo.Contains(memberRolesIds, liveroleId) {
				err = c.router.Session().GuildMemberRoleRemove(u.GuildID, u.User.ID, liveroleId)
				if err != nil {
					c.logger.Error("could not remove role", err)
				}
			}
		}
	}

	// c.logger.Debug("Status changed",
	// 	slog.String("user", u.Presence.User.ID),
	// 	slog.String("status", string(u.Status)),
	// )
}
