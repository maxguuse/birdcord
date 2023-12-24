package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrUserSide = errors.New("user side error")

	ErrWrongPollOptionLength  = errors.New("invalid option length")
	ErrWrongPollOptionsAmount = errors.New("invalid options amount")
	ErrAlreadyVoted           = errors.New("already voted")
	ErrNotAuthor              = errors.New("not author")
	ErrWrongGuild             = errors.New("wrong guild")
)
