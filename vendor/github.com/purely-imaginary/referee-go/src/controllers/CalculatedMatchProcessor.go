package main

import (
	"math"
	"time"

	"github.com/purely-imaginary/referee-go/src/models"
)

func processPlayersFromTeam(players []string, isRed bool) ([]PlayerSnapshot, float32) {
	var parsedPlayers []PlayerSnapshot
	var ratingSum float32 = 0
	for _, rawPlayerName := range players {
		player := GetPlayerByName(rawPlayerName)
		if player == nil {
			var newPlayer Player
			newPlayer.Name = rawPlayerName
			newPlayer.Wins = 0
			newPlayer.Losses = 0
			newPlayer.GoalsShot = 0
			newPlayer.GoalsScored = 0
			newPlayer.GoalsLost = 0
			newPlayer.WinRate = 0
			newPlayer.Rating = 1000
			newPlayer.Matches = []CalculatedMatch{}
			newPlayer.ID = newPlayer.InsertIntoDB()
			player = &newPlayer
		}
		parsedPlayers = append(parsedPlayers, *player.CreateSnapshot(isRed))
		ratingSum += player.Rating
	}
	return parsedPlayers, (ratingSum / float32(len(parsedPlayers)))
}

func processTime(timeString string) int64 {
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		panic("Can't parse time format")
	}

	return t.Unix()
}

func calculateRatingChange(calculatedMatch CalculatedMatch) float32 {
	kCoefficient := float32(250)

	ratingDifference := calculatedMatch.BlueTeam.AvgTeamRating - calculatedMatch.RedTeam.AvgTeamRating
	powerPiece := math.Pow(10, float64(ratingDifference/400))
	winChance := float32(1 / (1 + powerPiece))
	scorePerformance := float32(0.5)
	if !(calculatedMatch.RedTeam.Score+calculatedMatch.BlueTeam.Score == 0) {
		scoreDifference := calculatedMatch.RedTeam.Score - calculatedMatch.BlueTeam.Score
		if scoreDifference > 0 {
			scorePerformance = (((float32(1-winChance) / 10) * float32(scoreDifference)) + winChance)
		} else {
			scorePerformance = (((float32(winChance) / 10) * float32(scoreDifference)) + winChance)
		}
		// Old calc method:
		// scorePerformance = float32(scoreDifference+10) / 20
	}
	ratingChange := (scorePerformance - winChance) * kCoefficient
	ratingChangePerPlayer := ratingChange / float32(len(calculatedMatch.RedTeam.Players))
	return ratingChangePerPlayer
}

func updatePlayers(cm CalculatedMatch) {
	for _, snap := range cm.RedTeam.Players {
		UpdatePlayer(
			snap.PlayerID,
			cm.RedTeam.Score > cm.BlueTeam.Score,
			cm.RedTeam.Score,
			cm.BlueTeam.Score,
			cm.RedTeam.RatingChange)
		snap.UpdateMatchID(cm.ID)
	}
	for _, snap := range cm.BlueTeam.Players {
		UpdatePlayer(
			snap.PlayerID,
			cm.RedTeam.Score < cm.BlueTeam.Score,
			cm.BlueTeam.Score,
			cm.RedTeam.Score,
			cm.BlueTeam.RatingChange)
		snap.UpdateMatchID(cm.ID)
	}
}

func calculateMatch(rawMatch RawMatch) CalculatedMatch {
	duplicateID := CheckForDuplicatePositions(rawMatch.RawPositionsAtEnd)
	if duplicateID != 0 {
		return *GetMatchByID(duplicateID)
	}

	var processedMatch CalculatedMatch
	processedMatch.RawPositions = rawMatch.RawPositionsAtEnd

	processedMatch.RedTeam.Players, processedMatch.RedTeam.AvgTeamRating = processPlayersFromTeam(rawMatch.Teams.Red, true)
	processedMatch.BlueTeam.Players, processedMatch.BlueTeam.AvgTeamRating = processPlayersFromTeam(rawMatch.Teams.Blue, false)

	processedMatch.Time = rawMatch.Time
	processedMatch.RedTeam.Score = rawMatch.Score.Red
	processedMatch.BlueTeam.Score = rawMatch.Score.Blue

	processedMatch.RedTeam.RatingChange = calculateRatingChange(processedMatch)
	processedMatch.BlueTeam.RatingChange = -processedMatch.RedTeam.RatingChange

	processedMatch.ID = processedMatch.InsertToDB()
	updatePlayers(processedMatch)

	return processedMatch

}
