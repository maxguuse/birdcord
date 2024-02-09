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
)
