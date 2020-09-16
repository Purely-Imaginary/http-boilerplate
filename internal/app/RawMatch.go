package main

import (
	"github.com/jinzhu/gorm"
	"github.com/rushteam/gosql"
)

// RawGoal ..
type RawGoal struct {
	gorm.Model
	GoalScorerName string
	GoalShotTime   float32
	GoalSide       string
	GoalTime       float32
	GoalTravelTime float32
}

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
	GoalsData []RawGoal
}

// InsertIntoDB .
func (rm RawMatch) InsertIntoDB(DBEngine *gosql.PoolCluster) int64 {
	result, err := DBEngine.Insert(&rm)
	Check(err)
	insertID, _ := result.LastInsertId()
	return insertID
}
