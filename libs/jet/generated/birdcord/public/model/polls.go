//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Polls struct {
	ID        int32 `sql:"primary_key"`
	Title     string
	IsActive  bool
	CreatedAt *time.Time
	GuildID   int32
	AuthorID  *int32
}
