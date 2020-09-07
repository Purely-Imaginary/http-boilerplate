package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"

	"github.com/purely-imaginary/referee-go/src/models"
	"github.com/purely-imaginary/referee-go/src/tools"
)

//UnparsedReplayFolder comment
const UnparsedReplayFolder = "files/unparsedReplays/"

// ParsedReplayFolder comment
const ParsedReplayFolder = "files/replayData/"

// ParseReplay takes unparsed data and executes node converter to bin and then pythons exctractor of important data
func ParseReplay(replayName string) {
	filepath := UnparsedReplayFolder + replayName

	nodeCoverterPath := "hb-parser/haxball/replay.js"
	targetPath := "hb-parser/preprocessed/" + replayName + ".bin"

	cmd := exec.Command("node", nodeCoverterPath, "convert", filepath, targetPath)
	_, err := cmd.Output()
	Check(err)

	cmd = exec.Command("python3", "test.py", "preprocessed/"+replayName+".bin")
	cmd.Dir = "hb-parser"
	_, err = cmd.Output()
	Check(err)
	log.Println(replayName + " parsed")
}

// AsyncParseReplay ..
func AsyncParseReplay(replayName string, wg *sync.WaitGroup) {
	ParseReplay(replayName)
	wg.Done()
}

// ReadMatchFromFile takes converted file and creates RawMatch object from it
func ReadMatchFromFile(replayName string) RawMatch {
	suffix := ".bin.json"
	filepath := ParsedReplayFolder + replayName + suffix
	byteValue, err := ioutil.ReadFile(filepath)
	Check(err)

	var match RawMatch
	json.Unmarshal(byteValue, &match)

	return match
}
