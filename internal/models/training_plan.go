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

type TrainingPlanResponse struct {
	ID          uint   `json:"id"`
	TrainerID   uint   `json:"trainer-id"`
	ClientID    uint   `json:"client-id"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	CreatedAt   string `json:"created-at"`
}
