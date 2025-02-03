package models

import "time"

type TrainingPlan struct {
	ID          uint   `gorm:"primaryKey"`
	TrainerID   uint   `gorm:"not null"`
	ClientID    uint   `gorm:"not null"`
	Description string `gorm:"not null"`
	Schedule    string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
