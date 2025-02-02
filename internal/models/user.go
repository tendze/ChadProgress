package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"type:varchar(100);unique;not null"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Role         string    `gorm:"type:role_enum;not null"`
	RegisteredAt time.Time `gorm:"autoCreateTime"`
	Trainer      Trainer   `gorm:"foreignKey:UserID"`
	Client       Client    `gorm:"foreignKey:UserID"`
}
