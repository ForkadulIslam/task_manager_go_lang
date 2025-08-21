package models

import "time"

// Group represents the group model
type Group struct {
	ID        uint      `gorm:"primaryKey"`
	Label     string    `gorm:"not null;unique"`
	CreatedBy uint      `gorm:"not null"` // FK to users.id
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
