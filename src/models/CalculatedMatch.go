package models

import (
	"github.com/purely-imaginary/referee-go/src/controllers"
	"github.com/purely-imaginary/referee-go/src/tools"
)

// Team .
type Team struct {
	Players       []PlayerSnapshot
	AvgTeamRating float32
	Score         int64
	RatingChange  float32
}

//CalculatedMatch - match filled with full match and player data
type CalculatedMatch struct {
	ID           int64        `db:"id"`
	Time         string       `db:"time"`
	RedTeam      TeamSnapshot `gorm:"foreignkey:red_team_snapshot`
	BlueTeam     TeamSnapshot `gorm:"foreignkey:blue_team_snapshot`
	RawPositions string       `db:"raw_positions"`
	Goals        []Goal       `gorm:"foreignkey:goal_id`
}

//InsertToDB .
func (cm *CalculatedMatch) InsertToDB() int64 {
	err := controllers.DBEngine.Save(cm)
	tools.Check(err.Error)

	return cm.ID
}

// GetMatchByID .
func GetMatchByID(id int64) *CalculatedMatch {
	cm := &CalculatedMatch{}
	err := controllers.DBEngine.First(cm, "id = ?", id)

	if err.Error != nil {
		return nil
	}

	return cm
}

// CheckForDuplicatePositions .
func CheckForDuplicatePositions(positions string) int64 {
	cm := &CalculatedMatch{}
	err := controllers.DBEngine.First(cm, "raw_positions = ?", positions)

	if err.Error != nil {
		return 0
	}
	return cm.ID
}

// GetLastMatches ..
func GetLastMatches(amount int) []CalculatedMatch {
	var cms []CalculatedMatch
	err := controllers.DBEngine.Order("id DESC").Limit(amount).Find(&cms)

	if err.Error != nil {
		return nil
	}

	return cms
}

/*
type SQLPlayerHistory struct {
	ID       int64   `db:"uid,pk"`
	MatchID  int64   `db:"match_id"`
	PlayerID int64   `db:"player_id"`
	Time     int64   `db:"timestamp"`
	Rating   float32 `db:"rating"`
	IsRed    bool    `db:"is_red"`
}
*/
