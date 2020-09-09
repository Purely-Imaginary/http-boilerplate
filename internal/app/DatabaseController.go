package main

import (
	"github.com/jinzhu/gorm"

	//mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DBEngine .
var DBEngine *gorm.DB

//InitializeDB initializes DB if not already existing
func InitializeDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/referee?parseTime=true&readTimeout=1s&writeTimeout=1s&timeout=1s")
	Check(err)

	return db
}

// Migrate ..
func Migrate() {
	DBEngine.AutoMigrate(
		&CalculatedMatch{},
		&Player{},
		&PlayerSnapshot{},
		&DownloadedURL{},
	)
}

// DeleteAll .
func DeleteAll() {
	DBEngine.Exec("drop table referee.downloaded_url;")
	DBEngine.Exec("drop table referee.match_calculated;")
	DBEngine.Exec("drop table referee.player;")
	DBEngine.Exec("drop table referee.player_snapshot;")
	DBEngine.Exec("drop table referee.raw_match;")
}

// TruncateAll .
func TruncateAll() {
	DBEngine.Exec("truncate table referee.downloaded_url;")
	DBEngine.Exec("truncate table referee.match_calculated;")
	DBEngine.Exec("truncate table referee.player;")
	DBEngine.Exec("truncate table referee.player_snapshot;")
	DBEngine.Exec("truncate table referee.raw_match;")
}
