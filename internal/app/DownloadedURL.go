package main

import (
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// DownloadedURL ..
type DownloadedURL struct {
	URL     string `db:"url"`
	MatchID uint   `db:"match_id"`
}

// TableName .
func (u *DownloadedURL) TableName() string {
	return "downloaded_url"
}

// DoesExistsInDB .
func (u *DownloadedURL) DoesExistsInDB() uint {
	url := &DownloadedURL{}
	err := DBEngine.First(url, "url = ?", u.URL)

	if gorm.IsRecordNotFoundError(err.Error) {
		return 0
	}

	if url.MatchID == 0 {
		log.Println("Waiting for processing replay")
		for i := 0; i < 30; i++ {
			log.Println("W:" + strconv.Itoa(int(i)))
			time.Sleep(200 * time.Millisecond)
			checkURL := &DownloadedURL{}
			DBEngine.First(checkURL, "url = ?", u.URL)
			if checkURL.MatchID != 0 {
				log.Println("Found replay, ID:" + strconv.Itoa(int(checkURL.MatchID)))
				url.MatchID = checkURL.MatchID
				break
			}
		}
	}

	return url.MatchID
}

// InsertIntoDB .
func (u *DownloadedURL) InsertIntoDB() {
	DBEngine.Save(u)
}

// Update .
func (u *DownloadedURL) Update() {
	DBEngine.Model(&u).Update("match_id", u.MatchID)
	// url := DownloadedURL{}
	// DBEngine.First(&url, "url = ? AND match_id = 0", u.URL)
	// url.MatchID = u.MatchID
	// DBEngine.Save(&url)
}
