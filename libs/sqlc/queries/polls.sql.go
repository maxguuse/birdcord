// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: polls.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPoll = `-- name: CreatePoll :one
INSERT INTO polls (
    title, 
    author_id, 
    guild_id
) VALUES (
    $1, $2, $3
) RETURNING id, title, is_active, created_at, guild_id, author_id
`

type CreatePollParams struct {
	Title    string      `json:"title"`
	AuthorID pgtype.Int4 `json:"author_id"`
	GuildID  int32       `json:"guild_id"`
}

func (q *Queries) CreatePoll(ctx context.Context, arg CreatePollParams) (Poll, error) {
	row := q.db.QueryRow(ctx, createPoll, arg.Title, arg.AuthorID, arg.GuildID)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.IsActive,
		&i.CreatedAt,
		&i.GuildID,
		&i.AuthorID,
	)
	return i, err
}

const getActivePolls = `-- name: GetActivePolls :many
SELECT id, title, is_active, created_at, guild_id, author_id FROM polls WHERE is_active = true AND guild_id = $1 AND author_id = $2
`

type GetActivePollsParams struct {
	GuildID  int32       `json:"guild_id"`
	AuthorID pgtype.Int4 `json:"author_id"`
}

func (q *Queries) GetActivePolls(ctx context.Context, arg GetActivePollsParams) ([]Poll, error) {
	rows, err := q.db.Query(ctx, getActivePolls, arg.GuildID, arg.AuthorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Poll
	for rows.Next() {
		var i Poll
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.IsActive,
			&i.CreatedAt,
			&i.GuildID,
			&i.AuthorID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPoll = `-- name: GetPoll :one
SELECT id, title, is_active, created_at, guild_id, author_id FROM polls WHERE id = $1
`

func (q *Queries) GetPoll(ctx context.Context, id int32) (Poll, error) {
	row := q.db.QueryRow(ctx, getPoll, id)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.IsActive,
		&i.CreatedAt,
		&i.GuildID,
		&i.AuthorID,
	)
	return i, err
}

const updatePollStatus = `-- name: UpdatePollStatus :exec
UPDATE polls SET "is_active" = $2 WHERE id = $1
`

type UpdatePollStatusParams struct {
	ID       int32 `json:"id"`
	IsActive bool  `json:"is_active"`
}

func (q *Queries) UpdatePollStatus(ctx context.Context, arg UpdatePollStatusParams) error {
	_, err := q.db.Exec(ctx, updatePollStatus, arg.ID, arg.IsActive)
	return err
}
