package models

var (
	StatusBusy       = "BUSY"
	StatusActive     = "ACTIVE"
	StatusOnVacation = "ON_VACATION"
)

type Trainer struct {
	ID              uint             `gorm:"primaryKey"`
	UserID          uint             `gorm:"unique;not null"`
	Qualifications  string           `gorm:"type:varchar(150)"`
	Experience      string           `gorm:"type:varchar(250)"`
	Achievements    string           `gorm:"type:varchar(250)"`
	Status          string           `gorm:"type:status;default:'ACTIVE';not null"`
	Clients         []Client         `gorm:"foreignKey:TrainerID"`
	TrainingPlans   []TrainingPlan   `gorm:"foreignKey:TrainerID"`
	ProgressReports []ProgressReport `gorm:"foreignKey:TrainerID"`
}
