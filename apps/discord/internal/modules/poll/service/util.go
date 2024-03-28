package service

import (
	"errors"
	"strconv"
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

type voteData struct {
	PollId   int
	OptionId int
}

func parseVoteData(customId string) (*voteData, error) {
	parts := strings.Split(customId, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid custom id")
	}

	blob := parts[1]

	parts = strings.Split(blob, "_")
	if len(parts) != 6 {
		return nil, errors.New("invalid custom id")
	}

	poll_id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	option_id, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	return &voteData{
		PollId:   poll_id,
		OptionId: option_id,
	}, nil
}
