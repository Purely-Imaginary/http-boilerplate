package models

import (
	"github.com/purely-imaginary/referee-go/src/tools"
	"github.com/rushteam/gosql"
)

// RawMatch - data parsed from replay
type RawMatch struct {
	ID                int64  `db:"uid,pk"`
	Time              string `db:"time"`
	RawPositionsAtEnd string `db:"positions"`
	Teams             struct {
		Red  []string
		Blue []string
	}
	Score struct {
		Red  int64
		Blue int64
	}
}

// InsertIntoDB .
func (rm RawMatch) InsertIntoDB(DBEngine *gosql.PoolCluster) int64 {
	result, err := DBEngine.Insert(&rm)
	tools.Check(err)
	insertID, _ := result.LastInsertId()
	return insertID
}
