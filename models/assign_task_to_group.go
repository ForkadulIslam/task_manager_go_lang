package models

// AssignTaskToGroup represents the assignment of a task to a group
type AssignTaskToGroup struct {
	ID      uint `gorm:"primaryKey"`
	GroupID uint `gorm:"not null"` // FK to groups.id
	TaskID  uint `gorm:"not null"` // FK to tasks.id
}
