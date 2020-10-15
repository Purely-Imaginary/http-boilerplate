package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func prepareData(data interface{}, w http.ResponseWriter, r *http.Request) string {
	log.Println("Preparing data for ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	outputMessage, err := json.Marshal(&data)
	Check(err)
	return string(outputMessage)
}

// GetLastMatches ..
func GetLastMatches(w http.ResponseWriter, r *http.Request) {
	outputMessage := prepareData(GetLastMatchesFromDB(30), w, r)
	fmt.Fprintf(w, outputMessage)
}

// GetMatchByIDToAPI ..
func GetMatchByIDToAPI(w http.ResponseWriter, r *http.Request) {
	matchIDString := r.URL.Query()["id"][0]
	matchID, _ := strconv.Atoi(matchIDString)
	outputMessage := prepareData(GetMatchByID(uint(matchID)), w, r)
	fmt.Fprintf(w, outputMessage)
}

// GetPlayersTable ..
func GetPlayersTable(w http.ResponseWriter, r *http.Request) {
	outputMessage := prepareData(GetPlayersTableFromDB(), w, r)
	fmt.Fprintf(w, outputMessage)
}

// GetPlayersSnapshots ..
func GetPlayersSnapshots(w http.ResponseWriter, r *http.Request) {
	DBData := GetPlayersSnapshotsFromDB()
	for i := 0; i < len(DBData); i++ {
		snapshot := DBData[i]
		for _, targetSnapshot := range DBData {
			if snapshot.PlayerID == targetSnapshot.PlayerID &&
				snapshot.MatchRef.Time[:10] == targetSnapshot.MatchRef.Time[:10] &&
				snapshot.MatchRef.Time[11:] < targetSnapshot.MatchRef.Time[11:] {
				DBData = append(DBData[:i], DBData[i+1:]...)
				i--
				break
			}
		}
	}
	outputMessage := prepareData(DBData, w, r)
	fmt.Fprintf(w, outputMessage)
}

// GetFile ..
func GetFile(w http.ResponseWriter, r *http.Request) {
	matchIDString := r.URL.Query()["id"][0]
	matchID, _ := strconv.Atoi(matchIDString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fileURL := "files/unparsedReplays/" + GetMatchURLFromID(uint(matchID))
	http.ServeFile(w, r, fileURL)
}

// GetPlayerData ..
func GetPlayerData(w http.ResponseWriter, r *http.Request) {
	playerIDString := r.URL.Query()["id"][0]
	playerID, _ := strconv.Atoi(playerIDString)
	outputMessage := prepareData(GetPlayerDataFromDB(playerID), w, r)
	fmt.Fprintf(w, outputMessage)
}
