package models

import "time"

// Group represents the group model
type Group struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Label     string    `gorm:"not null;unique" json:"label"`
	CreatedBy uint      `gorm:"not null"` // FK to users.id
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime" json:"updated_at"`
	Users     []User    `gorm:"many2many:user_groups;joinForeignKey:group_id;joinReferences:user_id" json:"users"`
}

type GroupResponse struct {
	ID        uint      `json:"id"`
	Label     string    `json:"label"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Users     []UserResponse `json:"users"`
}