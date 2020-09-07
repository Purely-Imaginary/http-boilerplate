package models

// PlayerSnapshot - snapshot of player in history
type PlayerSnapshot struct {
	ID         int64              `db:"id"`
	PlayerID   int64              `db:"player_id"`
	PlayerName string             `db:"player_name"`
	MatchID    int64              `db:"match_id"`
	MatchRef   SQLCalculatedMatch `gorm:"foreignkey:match_id"`
	Rating     float32            `db:"rating"`
	IsRed      bool               `db:"is_red"`
}

// TableName .
func (u *PlayerSnapshot) TableName() string {
	return "player_snapshot"
}

// UpdateMatchID .
func (u *PlayerSnapshot) UpdateMatchID(matchID int64) {
	snap := &PlayerSnapshot{}
	err := DBEngine.First(snap, "id = ?", u.ID)
	Check(err.Error)
	snap.MatchID = matchID

	err = DBEngine.Save(snap)
	Check(err.Error)
}

// GetPlayersSnapshots ..
func GetPlayersSnapshots() []PlayerSnapshot {
	var playerSnaps []PlayerSnapshot
	err := DBEngine.Preload("MatchRef").Order("player_snapshot.match_id ASC").Find(&playerSnaps)

	if err.Error != nil {
		return nil
	}

	return playerSnaps
}
