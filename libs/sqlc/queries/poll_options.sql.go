// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: poll_options.sql

package queries

import (
	"context"
)

const createPollOption = `-- name: CreatePollOption :many
INSERT INTO poll_options (
    title, 
    poll_id
) VALUES (
    $1, $2
) RETURNING id, title, poll_id
`

type CreatePollOptionParams struct {
	Title  string `json:"title"`
	PollID int32  `json:"poll_id"`
}

func (q *Queries) CreatePollOption(ctx context.Context, arg CreatePollOptionParams) ([]PollOption, error) {
	rows, err := q.db.Query(ctx, createPollOption, arg.Title, arg.PollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PollOption
	for rows.Next() {
		var i PollOption
		if err := rows.Scan(&i.ID, &i.Title, &i.PollID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
