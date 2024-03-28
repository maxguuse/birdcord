package service

import "github.com/maxguuse/birdcord/apps/discord/internal/domain"

var (
	ErrNotAuthor = &domain.UsersideError{
		Msg: "Для изменения опроса нужно быть его автором.",
	}

	ErrNotFound = &domain.UsersideError{
		Msg: "Опроса не существует.",
	}

	ErrOptionHasVotes = &domain.UsersideError{
		Msg: "Невозможно удалить вариант опроса за который кто-то проголосовал.",
	}

	ErrTooFewOptions = &domain.UsersideError{
		Msg: "В опросе не может быть менее 2х вариантов.",
	}

	ErrTooManyOptions = &domain.UsersideError{
		Msg: "В опросе не может быть более 25 вариантов.",
	}
)