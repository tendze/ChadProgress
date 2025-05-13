package models

import "time"

type Metric struct {
	ID         uint `gorm:"primaryKey"`
	ClientID   uint `gorm:"not null"`
	Height     float64
	Weight     float64
	BodyFat    float64
	BMI        float64
	MeasuredAt time.Time `gorm:"autoCreateTime"`
}

type MetricResponse struct {
	ID         uint    `json:"id"`
	ClientID   uint    `json:"client-id"`
	Height     float64 `json:"height"`
	Weight     float64 `json:"weight"`
	BodyFat    float64 `json:"bodyfat"`
	BMI        float64 `json:"bmi"`
	MeasuredAt string  `json:"measured-at"`
}
