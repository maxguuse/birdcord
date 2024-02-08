package domain

import "errors"

type UsersideError struct {
	Msg string
}

func (e *UsersideError) Error() string {
	return e.Msg
}

var (
	ErrInternal = errors.New("internal error")
)
