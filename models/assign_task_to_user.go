package models

// AssignTaskToUser represents the assignment of a task to a user
type AssignTaskToUser struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"` // FK to users.id
	TaskID uint `gorm:"not null"` // FK to tasks.id
}
