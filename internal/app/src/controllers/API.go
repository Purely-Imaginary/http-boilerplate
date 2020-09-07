package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	controllers "purely-imaginary/referee-go/src/models"
)

func prepareData(data interface{}, w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	outputMessage, err := json.Marshal(&data)
	controllers.Check(err)
	return string(outputMessage)
}

// GetLastMatches ..
func GetLastMatches(w http.ResponseWriter, r *http.Request) {
	outputMessage := prepareData(GetLastMatches(200), w, r)
	fmt.Fprintf(w, outputMessage)
}

// GetPlayersTable ..
func GetPlayersTable(w http.ResponseWriter, r *http.Request) {
	outputMessage := prepareData(GetPlayersTable(), w, r)
	fmt.Fprintf(w, outputMessage)
}

// GetPlayersSnapshots ..
func GetPlayersSnapshots(w http.ResponseWriter, r *http.Request) {
	outputMessage := prepareData(GetPlayersSnapshots(), w, r)
	fmt.Fprintf(w, outputMessage)
}
