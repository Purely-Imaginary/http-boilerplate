package main

// PlayerToMatch ..
type PlayerToMatch struct {
	CalculatedMatchID uint
	PlayerID          uint
}

// TableName .
func (u *PlayerToMatch) TableName() string {
	return "player_match"
}
