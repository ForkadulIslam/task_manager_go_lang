package models

import "time"

// TaskStatusUpdateLog logs task status changes
type TaskStatusUpdateLog struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"` // FK to tasks.id
	UserID    uint      `gorm:"not null"` // FK to users.id
	Status    string    `gorm:"type:enum('Pending','In Progress','In Review','Completed');default:'Pending';comment:0=Pending,1=In Progress,2=In Review,3=Completed"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
