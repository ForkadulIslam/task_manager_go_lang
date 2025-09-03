package models

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	TaskID    uint      `json:"task_id"`
	Task      Task      `json:"task"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
