package repository

import "errors"

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return "could not find " + e.Resource
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrUserNotFound  = &NotFoundError{
		Resource: "user",
	}
	ErrPollNotFound = &NotFoundError{
		Resource: "poll",
	}
	ErrGuildNotFound = &NotFoundError{
		Resource: "guild",
	}
	ErrPollOptionNotFound = &NotFoundError{
		Resource: "poll option",
	}
	ErrLiveroleNotFound = &NotFoundError{
		Resource: "liverole",
	}
)
