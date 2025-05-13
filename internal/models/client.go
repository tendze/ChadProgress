package models

type Client struct {
	ID              uint `gorm:"primaryKey"`
	UserID          uint `gorm:"unique;not null"`
	TrainerID       uint `gorm:"not null"`
	Height          float64
	Weight          float64
	BodyFat         float64
	TrainingPlans   []TrainingPlan   `gorm:"foreignKey:ClientID"`
	ProgressReports []ProgressReport `gorm:"foreignKey:ClientID"`
	Metrics         []Metric         `gorm:"foreignKey:ClientID"`
}

type ClientResponse struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user-id"`
	TrainerID uint    `json:"trainer-id"`
	Height    float64 `json:"height"`
	Weight    float64 `json:"weight"`
	BodyFat   float64 `json:"bodyfat"`
}
