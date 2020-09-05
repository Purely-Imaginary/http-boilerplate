package models

import (
	"strings"

	"../repositories"
	"../tools"
	"github.com/jinzhu/gorm"
)

// Player - person who plays
type Player struct {
	ID          int64             `db:"id"`
	Name        string            `db:"name"`
	Wins        int64             `db:"wins"`
	Losses      int64             `db:"losses"`
	GoalsShot   int64             `db:"goals_shot"`
	GoalsScored int64             `db:"goals_scored"`
	GoalsLost   int64             `db:"goals_lost"`
	WinRate     float32           `db:"win_rate"`
	Rating      float32           `db:"rating"`
	Matches     []CalculatedMatch `gorm:"many2many:player_matches"`
}

// GetPlayerByName .
func GetPlayerByName(name string) *Player {
	player := &Player{}
	err := repositories.DBEngine.First(player, "name = ?", name)
	if gorm.IsRecordNotFoundError(err.Error) {
		return nil
	}
	tools.Check(err.Error)

	return getPlayerFromSQLPlayer(*SQLPlayer, true)
}

// GetPlayerByID .
func GetPlayerByID(id int) *Player {
	SQLPlayer := &repositories.SQLPlayer{}
	err := repositories.DBEngine.First(SQLPlayer, "id = ?", id)
	if gorm.IsRecordNotFoundError(err.Error) {
		return nil
	}
	tools.Check(err.Error)

	return getPlayerFromSQLPlayer(*SQLPlayer, true)
}

func getPlayerFromSQLPlayer(SQLPlayer repositories.SQLPlayer, withMatches bool) *Player {
	var playerSnapshots []PlayerSnapshot
	err := repositories.DBEngine.Find(&playerSnapshots, "player_id = ?", SQLPlayer.ID)
	tools.Check(err.Error)

	var matches []*CalculatedMatch
	var lastMatch string
	if withMatches {
		for _, playerMatch := range playerSnapshots {
			match := GetMatchByID(playerMatch.MatchID)
			matches = append(matches, match)
			lastMatch = match.Time
		}
	} else {
		match := GetMatchByID(playerSnapshots[SQLPlayer.Wins+SQLPlayer.Losses-1].MatchID)
		lastMatch = match.Time
	}

	lastMatch = lastMatch[:len(lastMatch)-1]
	lastMatch = strings.Replace(lastMatch, "T", " ", 1)

	returnObject := Player{
		ID:          SQLPlayer.ID,
		Name:        SQLPlayer.Name,
		Wins:        SQLPlayer.Wins,
		Losses:      SQLPlayer.Losses,
		GoalsScored: SQLPlayer.GoalsScored,
		GoalsLost:   SQLPlayer.GoalsLost,
		WinRate:     SQLPlayer.WinRate,
		Rating:      SQLPlayer.Rating,
		LastMatch:   lastMatch,
		Matches:     matches,
	}

	return &returnObject
}

// InsertIntoDB .
func (p *Player) InsertIntoDB() int64 {
	SQLObject := &repositories.SQLPlayer{
		Name:        p.Name,
		Wins:        p.Wins,
		Losses:      p.Losses,
		GoalsScored: p.GoalsScored,
		GoalsLost:   p.GoalsLost,
		WinRate:     p.WinRate,
		Rating:      p.Rating,
	}
	err := repositories.DBEngine.Save(SQLObject)
	tools.Check(err.Error)

	return SQLObject.ID
}

// CreateSnapshot .
func (p *Player) CreateSnapshot(isRed bool) *PlayerSnapshot {
	snapshot := &PlayerSnapshot{
		PlayerID:   p.ID,
		PlayerName: p.Name,
		Rating:     p.Rating,
		IsRed:      isRed,
	}
	err := repositories.DBEngine.Save(snapshot)
	tools.Check(err.Error)
	return snapshot
}

// UpdatePlayer .
func UpdatePlayer(PlayerID int64, win bool, goalsScored int64, goalsLost int64, ratingChange float32) {
	player := &repositories.SQLPlayer{}
	err := repositories.DBEngine.First(player, "id = ?", PlayerID)
	tools.Check(err.Error)
	if win {
		player.Wins = player.Wins + 1
	} else {
		player.Losses = player.Losses + 1
	}
	player.WinRate = float32(float32(player.Wins) / float32(player.Wins+player.Losses))
	player.GoalsScored += goalsScored
	player.GoalsLost += goalsLost
	player.Rating += ratingChange

	err = repositories.DBEngine.Save(player)
	tools.Check(err.Error)

}

// GetPlayersTable ..
func GetPlayersTable() []Player {
	var SQLObjects []repositories.SQLPlayer
	err := repositories.DBEngine.Order("rating DESC").Find(&SQLObjects)

	if err.Error != nil {
		return nil
	}
	var returnData []Player

	for _, SQLObject := range SQLObjects {
		returnPlayer := getPlayerFromSQLPlayer(SQLObject, false)
		returnData = append(returnData, *returnPlayer)
	}

	return returnData
}
