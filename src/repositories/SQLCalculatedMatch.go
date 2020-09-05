package repositories

// SQLCalculatedMatch .
type SQLCalculatedMatch struct {
	ID           int64   `db:"id"`
	Time         string  `db:"timestamp"`
	RedScore     int64   `db:"red_score"`
	BlueScore    int64   `db:"blue_score"`
	RedAvg       float32 `db:"red_avg"`
	BlueAvg      float32 `db:"blue_avg"`
	RatingChange float32 `db:"rating_change"`
	RawPositions string  `db:"raw_positions"`
}

// TableName .
func (u *SQLCalculatedMatch) TableName() string {
	return "match_calculated"
}

