//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Polls = newPollsTable("public", "polls", "")

type pollsTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	Title     postgres.ColumnString
	IsActive  postgres.ColumnBool
	CreatedAt postgres.ColumnTimestamp
	GuildID   postgres.ColumnInteger
	AuthorID  postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type PollsTable struct {
	pollsTable

	EXCLUDED pollsTable
}

// AS creates new PollsTable with assigned alias
func (a PollsTable) AS(alias string) *PollsTable {
	return newPollsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PollsTable with assigned schema name
func (a PollsTable) FromSchema(schemaName string) *PollsTable {
	return newPollsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PollsTable with assigned table prefix
func (a PollsTable) WithPrefix(prefix string) *PollsTable {
	return newPollsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PollsTable with assigned table suffix
func (a PollsTable) WithSuffix(suffix string) *PollsTable {
	return newPollsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPollsTable(schemaName, tableName, alias string) *PollsTable {
	return &PollsTable{
		pollsTable: newPollsTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newPollsTableImpl("", "excluded", ""),
	}
}

func newPollsTableImpl(schemaName, tableName, alias string) pollsTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		TitleColumn     = postgres.StringColumn("title")
		IsActiveColumn  = postgres.BoolColumn("is_active")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		GuildIDColumn   = postgres.IntegerColumn("guild_id")
		AuthorIDColumn  = postgres.IntegerColumn("author_id")
		allColumns      = postgres.ColumnList{IDColumn, TitleColumn, IsActiveColumn, CreatedAtColumn, GuildIDColumn, AuthorIDColumn}
		mutableColumns  = postgres.ColumnList{TitleColumn, IsActiveColumn, CreatedAtColumn, GuildIDColumn, AuthorIDColumn}
	)

	return pollsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		Title:     TitleColumn,
		IsActive:  IsActiveColumn,
		CreatedAt: CreatedAtColumn,
		GuildID:   GuildIDColumn,
		AuthorID:  AuthorIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
