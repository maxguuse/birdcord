package service

import (
	"strings"
	"unicode/utf8"

	"github.com/maxguuse/birdcord/apps/discord/internal/domain"
	"github.com/samber/lo"
)

func processPollOptions(rawOptions string) ([]string, error) {
	optionsList := strings.Split(rawOptions, "|")
	if len(optionsList) < 2 || len(optionsList) > 25 {
		return nil, &domain.UsersideError{
			Msg: "Количество вариантов опроса должно быть от 2 до 25 включительно.",
		}
	}
	if lo.SomeBy(optionsList, func(o string) bool {
		return utf8.RuneCountInString(o) > 50 || utf8.RuneCountInString(o) < 1
	}) {
		return nil, &domain.UsersideError{
			Msg: "Длина варианта опроса не может быть больше 50 или меньше 1 символа.",
		}
	}

	return optionsList, nil
}

func validatePollAuthor(poll *domain.PollWithDetails, userId string, guildId string) error {
	if poll.Author.DiscordUserID != userId {
		return ErrNotAuthor
	}

	if poll.Guild.DiscordGuildID != guildId {
		return ErrNotFound
	}

	return nil
}
