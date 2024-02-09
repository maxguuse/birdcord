package poll

import "github.com/maxguuse/birdcord/apps/discord/internal/domain"

var (
	ErrNotAuthor = &domain.UsersideError{
		Msg: "Для изменения опроса нужно быть его автором.",
	}

	ErrNotFound = &domain.UsersideError{
		Msg: "Опроса не существует.",
	}
)