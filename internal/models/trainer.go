package models

type Trainer struct {
	ID              uint `gorm:"primaryKey"`
	UserID          uint `gorm:"unique;not null"`
	Qualifications  string
	Experience      string
	Achievements    string
	Clients         []Client         `gorm:"foreignKey:TrainerID"`
	TrainingPlans   []TrainingPlan   `gorm:"foreignKey:TrainerID"`
	ProgressReports []ProgressReport `gorm:"foreignKey:TrainerID"`
}
