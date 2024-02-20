// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: liveroles.sql

package queries

import (
	"context"
)

const createLiveRole = `-- name: CreateLiveRole :one
INSERT INTO liveroles (guild_id, role_id) VALUES ($1, $2) RETURNING id, guild_id, role_id
`

type CreateLiveRoleParams struct {
	GuildID int32 `json:"guild_id"`
	RoleID  int32 `json:"role_id"`
}

func (q *Queries) CreateLiveRole(ctx context.Context, arg CreateLiveRoleParams) (Liverole, error) {
	row := q.db.QueryRow(ctx, createLiveRole, arg.GuildID, arg.RoleID)
	var i Liverole
	err := row.Scan(&i.ID, &i.GuildID, &i.RoleID)
	return i, err
}

const deleteLiveRoleByRoleID = `-- name: DeleteLiveRoleByRoleID :exec
DELETE FROM liveroles WHERE role_id = $1
`

func (q *Queries) DeleteLiveRoleByRoleID(ctx context.Context, roleID int32) error {
	_, err := q.db.Exec(ctx, deleteLiveRoleByRoleID, roleID)
	return err
}

const getLiveRolesByGuildID = `-- name: GetLiveRolesByGuildID :many
SELECT id, guild_id, role_id FROM liveroles WHERE guild_id = $1
`

func (q *Queries) GetLiveRolesByGuildID(ctx context.Context, guildID int32) ([]Liverole, error) {
	rows, err := q.db.Query(ctx, getLiveRolesByGuildID, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Liverole
	for rows.Next() {
		var i Liverole
		if err := rows.Scan(&i.ID, &i.GuildID, &i.RoleID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
