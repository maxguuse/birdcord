package repository

import (
	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"go.uber.org/fx"
)

var NewFx = fx.Options(
	fx.Provide(
		postgres.MustInit,

		NewGuildsRepository,
		NewPollsRepository,
		NewUsersRepository,

		fx.Annotate(
			NewDB,
			fx.As(new(DB)),
		),
	),
)

type DB interface {
	Polls() PollsRepository
	Users() UsersRepository
	Guilds() GuildsRepository
}

type db struct {
	pollsRepository  PollsRepository
	usersRepository  UsersRepository
	guildsRepository GuildsRepository
}

func NewDB(
	pr PollsRepository,
	ur UsersRepository,
	gr GuildsRepository,
) *db {
	return &db{
		pollsRepository:  pr,
		usersRepository:  ur,
		guildsRepository: gr,
	}
}

func (d *db) Polls() PollsRepository {
	return d.pollsRepository
}

func (d *db) Users() UsersRepository {
	return d.usersRepository
}

func (d *db) Guilds() GuildsRepository {
	return d.guildsRepository
}
