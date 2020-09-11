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
	DBEngine.AutoMigrate(&CalculatedMatch{})
	DBEngine.AutoMigrate(&Player{})
	DBEngine.AutoMigrate(&PlayerSnapshot{})
	DBEngine.AutoMigrate(&DownloadedURL{})
	DBEngine.AutoMigrate(&TeamSnapshot{})
	DBEngine.AutoMigrate(&Goal{})

}

// DeleteAll .
func DeleteAll() {
	DBEngine.Exec("drop table downloaded_url;")
	DBEngine.Exec("drop table calculated_matches;")
	DBEngine.Exec("drop table players;")
	DBEngine.Exec("drop table player_snapshot;")
	DBEngine.Exec("drop table goal;")
	DBEngine.Exec("drop table team_snapshot;")
}

// TruncateAll .
func TruncateAll() {
	DBEngine.Exec("truncate table downloaded_url;")
	DBEngine.Exec("truncate table calculated_matches;")
	DBEngine.Exec("truncate table players;")
	DBEngine.Exec("truncate table player_snapshot;")
}
