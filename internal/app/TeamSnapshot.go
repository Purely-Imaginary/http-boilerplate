package main

// TeamSnapshot ..
type TeamSnapshot struct {
	ID                    uint             `db:"id"`
	Players               []PlayerSnapshot `gorm:"foreignkey:team_snapshot_id"`
	AvgTeamRating         float32
	Score                 int64
	RedCalculatedMatchID  int64 `db:"red_calculated_match_id"`
	BlueCalculatedMatchID int64 `db:"blue_calculated_match_id"`
	RatingChange          float32
}

// TableName .
func (u *TeamSnapshot) TableName() string {
	return "team_snapshot"
}
