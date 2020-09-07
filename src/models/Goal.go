package models

// Goal ..
type Goal struct {
	ID         int64           `db:"id"`
	PlayerID   int64           `db:"player_id"`
	PlayerName string          `db:"player_name"`
	MatchID    int64           `db:"match_id"`
	MatchRef   CalculatedMatch `gorm:"foreignkey:match_id"`
	Time       float32         `db:"time"`
	IsRed      bool            `db:"is_red"`
}

// TableName .
func (u *Goal) TableName() string {
	return "goal"
}
