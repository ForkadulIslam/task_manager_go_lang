package models

import "time"

// TaskFollowupUser represents users assigned for task follow-up
type TaskFollowupUser struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"` // FK to users.id
	TaskID    uint      `gorm:"not null"` // FK to tasks.id
	Remarks   string    `gorm:"type:text;nullable"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserID"`
}
