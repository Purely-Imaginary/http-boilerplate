package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var startingRating float32 = 1000.0

func regenerateData(w http.ResponseWriter, r *http.Request) {
	fullRegenerate := false
	startTime := time.Now()
	DeleteAll()
	Migrate()
	if fullRegenerate {
		unparsedReplaysCounter := 1
		unparsedReplaysFiles, _ := ioutil.ReadDir(UnparsedReplayFolder)
		unparsedReplaysTotal := len(unparsedReplaysFiles)
		var wg sync.WaitGroup
		filepath.Walk(UnparsedReplayFolder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".hbr2" && strings.Contains(path, "HBReplay") {
				log.Print(strconv.Itoa(unparsedReplaysCounter) + "/" + strconv.Itoa(unparsedReplaysTotal))
				go AsyncParseReplay(info.Name(), &wg)
				unparsedReplaysCounter++
				wg.Add(1)
			}
			if unparsedReplaysCounter%5 == 0 {
				wg.Wait()
			}
			return nil
		})
		wg.Wait()
	}

	counter := 1
	files, _ := ioutil.ReadDir(ParsedReplayFolder)
	total := len(files)

	filepath.Walk(ParsedReplayFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".json" && strings.Contains(path, "HBReplay") {
			log.Print(info.Name() + " (" + strconv.Itoa(counter) + "/" + strconv.Itoa(total) + ")")
			ProcessReplayFromFile(strings.Trim(info.Name(), ".bin.json"))
			counter++
		}
		return nil
	})
	log.Println("\n" + time.Now().Sub(startTime).String())

}

func regenerateParsedReplays(w http.ResponseWriter, r *http.Request) {
	filepath.Walk(UnparsedReplayFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".hbr2" {
			ParseReplay(info.Name())
		}
		return nil
	})

}

func parseReplay(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	matchID := ProcessReplay(r)

	match := GetMatchByID(matchID)

	outputMessage := ExportHTML(*match)
	fmt.Fprintf(w, outputMessage+"\n"+(time.Now().Sub(startTime).String()))
}

func findTeams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	playerNames, _ := r.URL.Query()["players[]"]
	if len(playerNames) < 2 {
		fmt.Fprintf(w, string("not enough players"))
		return
	}
	var playerIDs []int
	playerData := make(map[int]*Player)
	for _, playerName := range playerNames {
		playerObject := GetPlayerByName(playerName)
		playerData[int(playerObject.ID)] = playerObject
		playerIDs = append(playerIDs, int(playerObject.ID))
	}
	allPermutations := Permutations(playerIDs)
	minDiff := 10000.00
	var permID int
	for permutationID, permutation := range allPermutations {
		var team1sum, team2sum float32
		for i := 0; i < len(permutation)/2; i++ {
			team1sum += playerData[permutation[i]].Rating
		}
		for i := len(permutation) / 2; i < len(permutation); i++ {
			team2sum += playerData[permutation[i]].Rating
		}
		diff := math.Abs(float64(team1sum - team2sum))
		if diff < minDiff {
			minDiff = diff
			permID = permutationID

		}
	}
	var teamRed, teamBlue []Player
	var redRatingSum, blueRatingSum = float32(0.0), float32(0.0)
	for i := 0; i < len(allPermutations[permID])/2; i++ {
		teamRed = append(teamRed, *playerData[allPermutations[permID][i]])
		redRatingSum += playerData[allPermutations[permID][i]].Rating
	}
	for i := len(allPermutations[permID]) / 2; i < len(allPermutations[permID]); i++ {
		teamBlue = append(teamBlue, *playerData[allPermutations[permID][i]])
		blueRatingSum += playerData[allPermutations[permID][i]].Rating
	}
	diff := float32(minDiff) / float32(len(teamRed))
	redRating := float64(float32(redRatingSum) / float32(len(teamRed)))
	blueRating := float64(float32(blueRatingSum) / float32(len(teamRed)))

	redChangePerGoal, blueChangePerGoal := CalculateFromRatings(float32(redRating), float32(blueRating), len(teamRed))

	outputData := struct {
		Red        []Player
		Blue       []Player
		Diff       float32
		BlueRating int
		RedRating  int
		BlueChange float32
		RedChange  float32
	}{
		teamRed,
		teamBlue,
		diff,
		int(math.Round(blueRating)),
		int(math.Round(redRating)),
		redChangePerGoal,
		blueChangePerGoal,
	}
	outputMessage, err := json.Marshal(&outputData)
	Check(err)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(outputMessage))
}

func main() {
	DBEngine = InitializeDB()
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/regenerate", regenerateData)
	http.HandleFunc("/getLastMatches", GetLastMatches)
	http.HandleFunc("/getPlayersTable", GetPlayersTable)
	http.HandleFunc("/getPlayersSnapshots", GetPlayersSnapshots)
	http.HandleFunc("/getMatchByID", GetMatchByIDToAPI)
	http.HandleFunc("/p", parseReplay)
	http.HandleFunc("/findTeams", findTeams)
	http.HandleFunc("/getFile", GetFile)
	http.HandleFunc("/getPlayerData", GetPlayerData)
	http.ListenAndServe(":7777", nil)
	log.Println("Ready to serve")
}

//HelloServer - testing function
func HelloServer(w http.ResponseWriter, r *http.Request) {
	InitializeDB()
	outputMessage := "<html><head><meta name=\"description\" content=\"" + "123456" + "\"></head><body>"
	outputMessage += "<form action=\"/regenerate\"><input type=\"submit\" value=\"Regenerate\" /></form>"
	outputMessage += "</body></html>"
	fmt.Fprintf(w, outputMessage)
}
