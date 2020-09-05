package repositories

// SQLPlayer .
type SQLPlayer struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Wins        int64   `db:"wins"`
	Losses      int64   `db:"losses"`
	GoalsScored int64   `db:"gwon"`
	GoalsLost   int64   `db:"glost"`
	WinRate     float32 `db:"winrate"`
	Rating      float32 `db:"current_rating"`
}

// TableName .
func (u *SQLPlayer) TableName() string {
	return "player"
}
