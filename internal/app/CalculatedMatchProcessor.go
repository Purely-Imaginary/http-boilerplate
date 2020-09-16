package main

import (
	"math"
	"time"
)

func processPlayersFromTeam(players []string, isRed bool) TeamSnapshot {
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
	var parsedTeam TeamSnapshot
	parsedTeam.Players = parsedPlayers
	parsedTeam.AvgTeamRating = ratingSum / float32(len(parsedPlayers))
	return parsedTeam
}

func processGoals(goals []RawGoal) []Goal {
	var parsedGoals []Goal
	for _, goal := range goals {
		player := &Player{}
		DBEngine.First(player, "name = ?", goal.GoalScorerName)

		var parsedGoal Goal
		parsedGoal.PlayerID = player.ID
		parsedGoal.PlayerName = goal.GoalScorerName
		parsedGoal.Time = goal.GoalTime
		parsedGoal.TravelTime = goal.GoalTravelTime
		parsedGoal.ShotTime = goal.GoalShotTime

		isRed := true
		if goal.GoalSide != "Red" {
			isRed = false
		}

		parsedGoal.IsRed = isRed
		parsedGoals = append(parsedGoals, parsedGoal)
	}
	return parsedGoals
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

	processedMatch.RedTeam = processPlayersFromTeam(rawMatch.Teams.Red, true)
	processedMatch.BlueTeam = processPlayersFromTeam(rawMatch.Teams.Blue, false)

	processedMatch.Goals = processGoals(rawMatch.GoalsData)
	processedMatch.Time = rawMatch.Time
	processedMatch.RedTeam.Score = rawMatch.Score.Red
	processedMatch.BlueTeam.Score = rawMatch.Score.Blue

	processedMatch.RedTeam.RatingChange = calculateRatingChange(processedMatch)
	processedMatch.BlueTeam.RatingChange = -processedMatch.RedTeam.RatingChange

	processedMatch.ID = processedMatch.InsertToDB()
	updatePlayers(processedMatch)

	return processedMatch

}
