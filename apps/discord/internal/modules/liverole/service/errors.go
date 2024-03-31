package service

import "github.com/maxguuse/birdcord/apps/discord/internal/domain"

var (
	ErrRoleAlreadyExists = &domain.UsersideError{
		Msg: "Данная роль уже добавлена.",
	}

	ErrNoLiveroles = &domain.UsersideError{
		Msg: "Нет live-ролей.",
	}

	ErrLiveroleNotFound = &domain.UsersideError{
		Msg: "Live-роль не найдена.",
	}
)
