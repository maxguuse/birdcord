package repository

import (
	"context"

	postgres "github.com/maxguuse/birdcord/libs/sqlc/db"
	"github.com/maxguuse/birdcord/libs/sqlc/queries"
)

type MessagesRepository interface {
	DeleteMessage(
		ctx context.Context,
		discordMessageId string,
	) error
}

type messagesRepository struct {
	q *postgres.DB
}

func NewMessagesRepository(q *postgres.DB) MessagesRepository {
	return &messagesRepository{
		q: q,
	}
}

func (m *messagesRepository) DeleteMessage(ctx context.Context, discordMessageId string) error {
	err := m.q.Transaction(ctx, func(q *queries.Queries) error {
		msg, err := q.GetMessageByDiscordID(ctx, discordMessageId)
		if err != nil {
			return err
		}

		err = q.DeleteMessageById(ctx, msg.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
