package models

import "time"

type ProgressReport struct {
	ID        uint `gorm:"primaryKey"`
	TrainerID uint `gorm:"not null"`
	ClientID  uint `gorm:"not null"`
	Comments  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
