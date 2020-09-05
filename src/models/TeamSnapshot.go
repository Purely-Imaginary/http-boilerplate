package models

import (
	"github.com/purely-imaginary/referee-go/src/models"
)

// TeamSnapshot ..
type TeamSnapshot struct {
	ID       int64                 `db:"id"`
	PlayerID models.PlayerSnapshot `gorm:"foreignkey:player_id"`
}

// TableName .
func (u *PlayerSnapshot) TableName() string {
	return "team_snapshot"
}
