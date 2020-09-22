package main

import "github.com/jinzhu/gorm"

// Team .
type Team struct {
	Players       []PlayerSnapshot
	AvgTeamRating float32
	Score         int64
	RatingChange  float32
}

//CalculatedMatch - match filled with full match and player data
type CalculatedMatch struct {
	gorm.Model
	Time         string       `db:"time"`
	StartTime    float32      `db:"start_time"`
	EndTime      float32      `db:"end_time"`
	RedTeam      TeamSnapshot `gorm:"foreignkey:red_calculated_match_id"`
	BlueTeam     TeamSnapshot `gorm:"foreignkey:blue_calculated_match_id"`
	RawPositions string       `gorm:"size:1000"`
	Goals        []Goal       `gorm:"foreignkey:match_id"`
}

//InsertToDB .
func (cm *CalculatedMatch) InsertToDB() uint {
	err := DBEngine.Save(cm)
	Check(err.Error)

	for _, goal := range cm.Goals {
		goal.MatchID = cm.ID
		goal.MatchRef = *cm

		DBEngine.Save(&goal)

		player := &Player{}
		DBEngine.First(player, "id = ?", goal.PlayerID)
		player.GoalsShot++
		DBEngine.Save(&player)
	}
	DBEngine.Save(&cm.RedTeam)
	DBEngine.Save(&cm.BlueTeam)

	for _, playerSnapshot := range append(cm.RedTeam.Players, cm.BlueTeam.Players...) {
		player := &Player{}
		DBEngine.First(player, "id = ?", playerSnapshot.PlayerID)
		var playerMatch PlayerToMatch
		playerMatch.PlayerID = (*player).ID
		playerMatch.CalculatedMatchID = cm.ID
		DBEngine.Save(&playerMatch)
	}
	return cm.ID
}

// GetMatchByID .
func GetMatchByID(id uint) *CalculatedMatch {
	cm := &CalculatedMatch{}
	err := DBEngine.
		Preload("Goals").
		Preload("RedTeam").
		Preload("BlueTeam").
		Preload("RedTeam.Players").
		Preload("BlueTeam.Players").
		First(cm, "id = ?", id)

	if err.Error != nil {
		return nil
	}

	return cm
}

// CheckForDuplicatePositions .
func CheckForDuplicatePositions(positions string) uint {
	cm := &CalculatedMatch{}
	err := DBEngine.First(cm, "raw_positions = ?", positions)

	if err.Error != nil {
		return 0
	}
	return cm.ID
}

// GetLastMatchesFromDB ..
func GetLastMatchesFromDB(amount int) []CalculatedMatch {
	var cms []CalculatedMatch
	err := DBEngine.Order("id DESC").
		Limit(amount).
		Preload("Goals").
		Preload("RedTeam").
		Preload("BlueTeam").
		Preload("RedTeam.Players").
		Preload("BlueTeam.Players").
		Find(&cms)

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
