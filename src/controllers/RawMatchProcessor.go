package controllers

import (
	"net/http"

	"../repositories"
)

func ProcessReplay(request *http.Request) int64 {
	downloadPath := ExtractURL(request)

	urlToCheck := repositories.SQLDownloadedUrl{URL: downloadPath}
	id := urlToCheck.DoesExistsInDB()
	if id != 0 {
		return id
	}

	replayName := DownloadReplay(downloadPath)

	rawMatch := ReadMatchFromFile(replayName)

	calculatedMatch := calculateMatch(rawMatch)

	urlToCheck.MatchID = calculatedMatch.ID
	urlToCheck.InsertIntoDB()

	return calculatedMatch.ID

}
