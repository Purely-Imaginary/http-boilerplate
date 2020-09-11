package main

// TeamSnapshot ..
type TeamSnapshot struct {
	ID            uint             `db:"id"`
	Players       []PlayerSnapshot `gorm:"foreignkey:player_id"`
	AvgTeamRating float32
	Score         int64
	RatingChange  float32
	IsRed         bool
}

// TableName .
func (u *TeamSnapshot) TableName() string {
	return "team_snapshot"
}
