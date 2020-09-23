package main

// Goal ..
type Goal struct {
	ID         uint            `db:"id"`
	PlayerID   uint            `db:"player_id"`
	PlayerName string          `db:"player_name"`
	MatchID    uint            `db:"match_id"`
	MatchRef   CalculatedMatch `gorm:"foreignkey:match_id"`
	Time       float32         `db:"time"`
	TravelTime float32         `db:"travel_time"`
	Speed      float32         `db:"speed"`
	ShotTime   float32         `db:"shot_time"`
	IsRed      bool            `db:"is_red"`
}

// TableName .
func (u *Goal) TableName() string {
	return "goal"
}
