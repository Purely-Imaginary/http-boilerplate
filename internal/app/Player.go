package main

import (
	"github.com/jinzhu/gorm"
)

// Player - person who plays
type Player struct {
	ID          uint              `db:"id"`
	Name        string            `db:"name"`
	Wins        int64             `db:"wins"`
	Losses      int64             `db:"losses"`
	GoalsShot   int64             `db:"goals_shot"`
	GoalsScored int64             `db:"goals_scored"`
	GoalsLost   int64             `db:"goals_lost"`
	WinRate     float32           `db:"win_rate"`
	Rating      float32           `db:"rating"`
	Matches     []CalculatedMatch `gorm:"many2many:player_match"`
}

// PlayerDetailedData - detailed data to show
type PlayerDetailedData struct {
	Player        Player
	Snapshots     []PlayerSnapshot
	PlayerRatings []Player
}

// GetPlayerByName .
func GetPlayerByName(name string) *Player {
	player := &Player{}
	err := DBEngine.First(player, "name = ?", name)
	if gorm.IsRecordNotFoundError(err.Error) {
		return nil
	}
	Check(err.Error)

	return player
}

// GetPlayerByID .
func GetPlayerByID(id int) *Player {
	player := &Player{}
	err := DBEngine.Preload("Matches").
		Preload("Matches.BlueTeam").
		Preload("Matches.BlueTeam.Players").
		Preload("Matches.RedTeam").
		Preload("Matches.RedTeam.Players").
		First(player, "id = ?", id)

	if gorm.IsRecordNotFoundError(err.Error) {
		return nil
	}
	Check(err.Error)

	return player
}

// InsertIntoDB .
func (p *Player) InsertIntoDB() uint {
	err := DBEngine.Save(p)
	Check(err.Error)

	return p.ID
}

// CreateSnapshot .
func (p *Player) CreateSnapshot(isRed bool) *PlayerSnapshot {
	snapshot := &PlayerSnapshot{
		PlayerID:   p.ID,
		PlayerName: p.Name,
		Rating:     p.Rating,
		IsRed:      isRed,
	}
	err := DBEngine.Save(snapshot)
	Check(err.Error)
	return snapshot
}

// UpdatePlayer .
func UpdatePlayer(PlayerID uint, win bool, goalsScored int64, goalsLost int64, ratingChange float32) {
	player := &Player{}
	err := DBEngine.First(player, "id = ?", PlayerID)
	if player.Name == "" {
		player = nil
	}
	Check(err.Error)
	if win {
		player.Wins = player.Wins + 1
	} else {
		player.Losses = player.Losses + 1
	}
	player.WinRate = float32(float32(player.Wins) / float32(player.Wins+player.Losses))
	player.GoalsScored += goalsScored
	player.GoalsLost += goalsLost
	player.Rating += ratingChange

	err = DBEngine.Save(player)
	Check(err.Error)

}

// GetPlayersTableFromDB ..
func GetPlayersTableFromDB() []Player {
	var players []Player
	err := DBEngine.Order("rating DESC").
		Preload("Matches").
		Where("wins + losses > 10").
		Find(&players)

	if err.Error != nil {
		return nil
	}

	return players
}

// GetPlayerSnapshots ..
func GetPlayerSnapshots(id int) []PlayerSnapshot {
	var playerSnaps []PlayerSnapshot
	err := DBEngine.Order("match_id ASC").
		Where("player_id = ?", id).
		Preload("MatchRef").
		Find(&playerSnaps)

	if err.Error != nil {
		return nil
	}

	return playerSnaps
}

// GetPlayerRatings ..
func GetPlayerRatings() []Player {
	var players []Player
	err := DBEngine.Order("rating DESC").
		Where("wins + losses > 10").
		Find(&players)

	if err.Error != nil {
		return nil
	}

	return players
}

// GetPlayerDataFromDB .
func GetPlayerDataFromDB(id int) *PlayerDetailedData {
	var returnValue PlayerDetailedData

	returnValue.Player = *GetPlayerByID(id)
	returnValue.Snapshots = GetPlayerSnapshots(id)
	returnValue.PlayerRatings = GetPlayerRatings()
	return &returnValue
}
