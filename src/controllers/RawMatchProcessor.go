package controllers

import (
	"net/http"
	"referee-go/src/models"
)

// ProcessReplay .
func ProcessReplay(request *http.Request) int64 {
	downloadPath := ExtractURL(request)

	urlToCheck := models.DownloadedURL{URL: downloadPath}
	id := urlToCheck.DoesExistsInDB()
	if id != 0 {
		return id
	}
	replayName := DownloadReplay(downloadPath)
	id = ProcessReplayFromFile(replayName)

	urlToCheck.MatchID = id
	urlToCheck.InsertIntoDB()

	return id

}

// ProcessReplayFromFile .
func ProcessReplayFromFile(replayName string) int64 {

	rawMatch := ReadMatchFromFile(replayName)

	calculatedMatch := calculateMatch(rawMatch)

	return calculatedMatch.ID
}
