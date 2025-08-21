package models

import "time"

// TaskType represents the task type model
type TaskType struct {
	ID        uint      `gorm:"primaryKey"`
	Label     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
