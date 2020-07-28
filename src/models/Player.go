package models

import (
	"../repositories"
	"../tools"
	"github.com/rushteam/gosql"
)

// Player - person who plays
type Player struct {
	ID          int64
	Name        string
	Wins        int64
	Losses      int64
	GoalsScored int64
	GoalsLost   int64
	WinRate     float32
	Rating      float32
	Matches     []RawMatch
}

func GetPlayerByName(name string) *Player {
	SQLPlayer := &repositories.SQLPlayer{}
	err := repositories.DBEngine.Fetch(SQLPlayer,
		gosql.Where("name", name),
	)

	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil
	}
	tools.Check(err)

	MatchesIDs := &repositories.SQLPlayerHistory{}
	err = repositories.DBEngine.FetchAll(MatchesIDs,
		gosql.Where("player_id", SQLPlayer.ID),
	)
	tools.Check(err)
	for _, MatchID := range MatchesIDs {

	}
	return url.MatchID
}
