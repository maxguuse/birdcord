//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type PollVotes struct {
	ID       int32 `sql:"primary_key"`
	PollID   int32
	OptionID int32
	UserID   int32
}
