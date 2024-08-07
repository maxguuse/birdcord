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

var PollVotes = newPollVotesTable("public", "poll_votes", "")

type pollVotesTable struct {
	postgres.Table

	// Columns
	ID       postgres.ColumnInteger
	PollID   postgres.ColumnInteger
	OptionID postgres.ColumnInteger
	UserID   postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type PollVotesTable struct {
	pollVotesTable

	EXCLUDED pollVotesTable
}

// AS creates new PollVotesTable with assigned alias
func (a PollVotesTable) AS(alias string) *PollVotesTable {
	return newPollVotesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PollVotesTable with assigned schema name
func (a PollVotesTable) FromSchema(schemaName string) *PollVotesTable {
	return newPollVotesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PollVotesTable with assigned table prefix
func (a PollVotesTable) WithPrefix(prefix string) *PollVotesTable {
	return newPollVotesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PollVotesTable with assigned table suffix
func (a PollVotesTable) WithSuffix(suffix string) *PollVotesTable {
	return newPollVotesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPollVotesTable(schemaName, tableName, alias string) *PollVotesTable {
	return &PollVotesTable{
		pollVotesTable: newPollVotesTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newPollVotesTableImpl("", "excluded", ""),
	}
}

func newPollVotesTableImpl(schemaName, tableName, alias string) pollVotesTable {
	var (
		IDColumn       = postgres.IntegerColumn("id")
		PollIDColumn   = postgres.IntegerColumn("poll_id")
		OptionIDColumn = postgres.IntegerColumn("option_id")
		UserIDColumn   = postgres.IntegerColumn("user_id")
		allColumns     = postgres.ColumnList{IDColumn, PollIDColumn, OptionIDColumn, UserIDColumn}
		mutableColumns = postgres.ColumnList{PollIDColumn, OptionIDColumn, UserIDColumn}
	)

	return pollVotesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:       IDColumn,
		PollID:   PollIDColumn,
		OptionID: OptionIDColumn,
		UserID:   UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
