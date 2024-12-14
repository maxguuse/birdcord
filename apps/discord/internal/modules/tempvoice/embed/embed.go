package embed

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

var (
	hubsManagementEmbed = []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Создание и настройка хабов",
			Description: "Выберите `Создать хаб`, чтобы начать процесс настройки нового хаба.\nВыберите существующий хаб чтобы изменить его конфигурацию.",
			Color:       0x5865f2,
		},
	}
	HubsManagement = func(
		channels []*discordgo.Channel,
	) ([]*discordgo.MessageEmbed, []discordgo.MessageComponent) {
		components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Создать хаб",
						Style: discordgo.SuccessButton,
						Emoji: &discordgo.ComponentEmoji{
							Name:     "➕",
							Animated: false,
						},
						CustomID: "create-tempvoice-hub-btn",
					},
				},
			},
		}

		if len(channels) > 0 {
			components = append(components, discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						MenuType:    discordgo.StringSelectMenu,
						CustomID:    "configure-tempvoice-hub-select-menu",
						Placeholder: "🔧 Изменить существующий хаб",
						Options: lo.Map(channels, func(c *discordgo.Channel, i int) discordgo.SelectMenuOption {
							return discordgo.SelectMenuOption{
								Label:       fmt.Sprintf("Изменить хаб #%d", i+1),
								Value:       c.ID,
								Description: c.Name,
								Emoji: &discordgo.ComponentEmoji{
									Name:     "🔧",
									Animated: false,
								},
								Default: false,
							}
						}),
						Disabled: false,
					},
				},
			})
		}

		return hubsManagementEmbed, components
	}
)

var (
	HubConfiguration = func(
		hub *domain.TempvoiceHub,
		hubName string,
	) ([]*discordgo.MessageEmbed, []discordgo.MessageComponent) {
		embed := []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Title:       "Настройка хаба",
				Description: "Пользователи могут присоединяться к хабу, чтобы создать временный голосовой канал.",
				Color:       0x5865f2,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "• Название хаба",
						Value:  hubName,
						Inline: false,
					},
					{
						Name:   "• Шаблон названия временных каналов",
						Value:  hub.TempvoiceTemplate,
						Inline: false,
					},
					{
						Name:   "• Категория для временных каналов",
						Value:  fmt.Sprintf("<#%d>", hub.TempvoiceCategory),
						Inline: false,
					},
				},
			},
		}

		components := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Изменить название",
						Style: discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							Name:     "✏",
							Animated: false,
						},
						CustomID: "", // TODO: Fill custom id
					},
					discordgo.Button{
						Label: "Изменить шаблон",
						Style: discordgo.SecondaryButton,
						Emoji: &discordgo.ComponentEmoji{
							Name:     "📐",
							Animated: false,
						},
						CustomID: "", // TODO: Fill custom id
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						MenuType:    discordgo.ChannelSelectMenu,
						CustomID:    "", // TODO: Fill custom id
						Placeholder: "📚 Выберите категорию",
						Disabled:    false,
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildCategory,
						},
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Подтвердить",
						Style:    discordgo.SuccessButton,
						CustomID: "", // TODO: Fill custom id
					},
					discordgo.Button{
						Label:    "Отменить",
						Style:    discordgo.DangerButton,
						CustomID: "", // TODO: Fill custom id
					},
				},
			},
		}

		return embed, components
	}
)
