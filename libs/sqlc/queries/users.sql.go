// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: users.sql

package queries

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    discord_user_id
) VALUES (
    $1
) ON CONFLICT (discord_user_id) DO NOTHING RETURNING id, discord_user_id
`

func (q *Queries) CreateUser(ctx context.Context, discordUserID string) (User, error) {
	row := q.db.QueryRow(ctx, createUser, discordUserID)
	var i User
	err := row.Scan(&i.ID, &i.DiscordUserID)
	return i, err
}

const getUserByDiscordID = `-- name: GetUserByDiscordID :one
SELECT id, discord_user_id FROM users WHERE discord_user_id = $1
`

func (q *Queries) GetUserByDiscordID(ctx context.Context, discordUserID string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByDiscordID, discordUserID)
	var i User
	err := row.Scan(&i.ID, &i.DiscordUserID)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, discord_user_id FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(&i.ID, &i.DiscordUserID)
	return i, err
}
