package models

import "time"

// TaskCommentLog logs comments on tasks
type TaskCommentLog struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"` // FK to tasks.id
	UserID    uint      `gorm:"not null"` // FK to users.id
	Comment   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
