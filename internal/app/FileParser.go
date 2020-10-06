package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
	"time"
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

	nodeTimeStart := time.Now()
	cmd := exec.Command("node", nodeCoverterPath, "convert", filepath, targetPath)
	_, err := cmd.Output()
	Check(err)
	nodeTimeElapsed := time.Now().Sub(nodeTimeStart).String()

	pythonTimeStart := time.Now()
	cmd = exec.Command("python3", "test.py", "preprocessed/"+replayName+".bin")
	cmd.Dir = "hb-parser"
	_, err = cmd.Output()

	if err != nil {
		log.Println("Python parser error: ", err)
	}

	pythonTimeElapsed := time.Now().Sub(pythonTimeStart).String()

	log.Println(replayName + " parsed. Node: " + nodeTimeElapsed + " - Python: " + pythonTimeElapsed)
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
