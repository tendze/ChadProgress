package models

import "time"

type Metric struct {
	ID         uint `gorm:"primaryKey"`
	ClientID   uint `gorm:"not null"`
	Weight     float64
	BodyFat    float64
	BMI        float64
	MeasuredAt time.Time `gorm:"autoCreateTime"`
}
