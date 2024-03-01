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
		NewLiverolesRepository,
		NewMessagesRepository,

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
	Liveroles() LiverolesRepository
	Messages() MessagesRepository
}

type db struct {
	pollsRepository     PollsRepository
	usersRepository     UsersRepository
	guildsRepository    GuildsRepository
	liverolesRepository LiverolesRepository
	messagesRepository  MessagesRepository
}

func NewDB(
	pr PollsRepository,
	ur UsersRepository,
	gr GuildsRepository,
	lr LiverolesRepository,
	mr MessagesRepository,
) *db {
	return &db{
		pollsRepository:     pr,
		usersRepository:     ur,
		guildsRepository:    gr,
		liverolesRepository: lr,
		messagesRepository:  mr,
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

func (d *db) Liveroles() LiverolesRepository {
	return d.liverolesRepository
}

func (d *db) Messages() MessagesRepository {
	return d.messagesRepository
}
