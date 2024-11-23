package service

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/maxguuse/birdcord/apps/discord/internal/modules/tempvoice/repository"
	"github.com/maxguuse/birdcord/libs/jet/txmanager"
	"github.com/samber/lo"
)

type Service struct {
	repo      repository.Repository
	txManager *txmanager.TxManager
}

func New(
	repo repository.Repository,
	txManager *txmanager.TxManager,
) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
	}
}

type GetHubsRequest struct {
	DiscordGuildId string
	Session        *discordgo.Session
	Interaction    *discordgo.Interaction
}

func (s *Service) GetHubs(
	ctx context.Context,
	req GetHubsRequest,
) []*discordgo.Channel { // TODO: Add error handling
	guildId, err := strconv.Atoi(req.DiscordGuildId)
	if err != nil {
		return nil
	}

	hubs, err := s.repo.GetHubs(ctx, int64(guildId))
	if err != nil {
		return nil
	}

	channels, err := req.Session.GuildChannels(req.DiscordGuildId)
	if err != nil {
		return nil
	}

	hubsIds := lo.Map(hubs, func(item *domain.TempvoiceHub, index int) int {
		return item.DiscordChannelID
	})

	filteredChannels := lo.Filter(channels, func(channel *discordgo.Channel, _ int) bool {
		channelId, _ := strconv.Atoi(channel.ID)
		return lo.Contains(hubsIds, channelId)
	})

	return filteredChannels
}
