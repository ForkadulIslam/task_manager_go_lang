package models

import "time"

// TaskStatusUpdateLog logs task status changes
type TaskStatusUpdateLog struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"` // FK to tasks.id
	UserID    uint      `gorm:"not null"` // FK to users.id
	Status    int       `gorm:"type:enum('0','1','2','3');comment:0=Pending,1=In Progress,2=In Review,3=Completed"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
