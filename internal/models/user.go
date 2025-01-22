package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Role         string    `gorm:"type:enum('trainer', 'client');not null"`
	RegisteredAt time.Time `gorm:"autoCreateTime"`
	Trainer      Trainer   `gorm:"foreignKey:UserID"`
	Client       Client    `gorm:"foreignKey:UserID"`
}
