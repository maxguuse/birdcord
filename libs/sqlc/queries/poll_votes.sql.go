// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: poll_votes.sql

package queries

import (
	"context"
)

const addVote = `-- name: AddVote :exec
INSERT INTO poll_votes (
    user_id, 
    poll_id, 
    option_id
) VALUES (
    $1, $2, $3
)
`

type AddVoteParams struct {
	UserID   int32 `json:"user_id"`
	PollID   int32 `json:"poll_id"`
	OptionID int32 `json:"option_id"`
}

func (q *Queries) AddVote(ctx context.Context, arg AddVoteParams) error {
	_, err := q.db.Exec(ctx, addVote, arg.UserID, arg.PollID, arg.OptionID)
	return err
}

const getPollVotes = `-- name: GetPollVotes :many
SELECT id, poll_id, option_id, user_id FROM poll_votes
WHERE poll_id = $1
`

func (q *Queries) GetPollVotes(ctx context.Context, pollID int32) ([]PollVote, error) {
	rows, err := q.db.Query(ctx, getPollVotes, pollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PollVote
	for rows.Next() {
		var i PollVote
		if err := rows.Scan(
			&i.ID,
			&i.PollID,
			&i.OptionID,
			&i.UserID,
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

const getVote = `-- name: GetVote :one
SELECT COUNT(*) FROM poll_votes 
WHERE user_id = $1 AND poll_id = $2
`

type GetVoteParams struct {
	UserID int32 `json:"user_id"`
	PollID int32 `json:"poll_id"`
}

func (q *Queries) GetVote(ctx context.Context, arg GetVoteParams) (int64, error) {
	row := q.db.QueryRow(ctx, getVote, arg.UserID, arg.PollID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
