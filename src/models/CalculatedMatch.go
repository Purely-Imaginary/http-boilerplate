package models

import (
	"../repositories"
)

//CalculatedMatch - match filled with full match and player data
type CalculatedMatch struct {
	ID    int64
	Time  int64
	Teams struct {
		Red struct {
			Players       []PlayerHistorical
			AvgTeamRating float32
			Score         int64
			RatingChange  float32
		}
		Blue struct {
			Players       []PlayerHistorical
			AvgTeamRating float32
			Score         int64
			RatingChange  float32
		}
	}
}

func (cm *CalculatedMatch) InsertToDB() int64 {
	match := repositories.SQLCalculatedMatch{
		Time:         cm.Time,
		RedScore:     cm.Teams.Red.Score,
		BlueScore:    cm.Teams.Blue.Score,
		RatingChange: cm.Teams.Red.RatingChange,
	}

	res, _ := repositories.DBEngine.Insert(match)
	matchID, _ := res.LastInsertId()
	for _, redPlayer := range cm.Teams.Red.Players {
		playerHistory := repositories.SQLPlayerHistory{
			MatchID:  matchID,
			PlayerID: redPlayer.ID,
			Rating:   redPlayer.Rating,
			IsRed:    true,
		}
		repositories.DBEngine.Insert(playerHistory)
	}
	for _, bluePlayer := range cm.Teams.Blue.Players {
		playerHistory := repositories.SQLPlayerHistory{
			MatchID:  matchID,
			PlayerID: bluePlayer.ID,
			Rating:   bluePlayer.Rating,
			IsRed:    false,
		}
		repositories.DBEngine.Insert(playerHistory)
	}
	return matchID
}

func GetMatchByID(id int64) {

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
