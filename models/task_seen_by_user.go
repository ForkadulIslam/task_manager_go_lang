package models

import "time"

// TaskSeenByUser logs when a task is seen by a user
type TaskSeenByUser struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"` // FK to tasks.id
	UserID    uint      `gorm:"not null"` // FK to users.id
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
