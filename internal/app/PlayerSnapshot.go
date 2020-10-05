package main

// PlayerSnapshot - snapshot of player in history
type PlayerSnapshot struct {
	ID             uint            `db:"id"`
	PlayerID       uint            `db:"player_id"`
	PlayerName     string          `db:"player_name"`
	TeamSnapshotID uint            `db:"team_snapshot_id"`
	MatchID        uint            `db:"match_id"`
	MatchRef       CalculatedMatch `gorm:"foreignkey:match_id"`
	Rating         float32         `db:"rating"`
	IsRed          bool            `db:"is_red"`
}

// TableName .
func (u *PlayerSnapshot) TableName() string {
	return "player_snapshot"
}

// UpdateMatchID .
func (u *PlayerSnapshot) UpdateMatchID(matchID uint) {
	snap := &PlayerSnapshot{}
	err := DBEngine.First(snap, "id = ?", u.ID)
	Check(err.Error)
	snap.MatchID = matchID

	err = DBEngine.Save(snap)
	Check(err.Error)
}

// GetPlayersSnapshotsFromDB ..
func GetPlayersSnapshotsFromDB() []PlayerSnapshot {
	var playerSnaps []PlayerSnapshot
	err := DBEngine.Preload("MatchRef").
		Preload("MatchRef.RedTeam").
		Order("player_snapshot.match_id ASC").
		Find(&playerSnaps)

	if err.Error != nil {
		return nil
	}

	for key := range playerSnaps {
		if playerSnaps[key].IsRed {
			playerSnaps[key].Rating = playerSnaps[key].Rating + playerSnaps[key].MatchRef.RedTeam.RatingChange
		} else {
			playerSnaps[key].Rating = playerSnaps[key].Rating - playerSnaps[key].MatchRef.RedTeam.RatingChange
		}
	}

	return playerSnaps
}
