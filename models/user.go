package models

import "time"

// User represents the user model
// @Description User model for the application
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // Exclude password from JSON
	Status    int       `gorm:"default:1;comment:1=Active,0=Inactive" json:"status"`
	UserLabel int       `gorm:"comment:1=Super Admin,2=User" json:"user_label"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime" json:"updated_at"`
}

type UserResponse struct {
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	Status        int    `json:"status"`
	UserLabel     int    `json:"user_label"`
	AssociationID uint   `json:"association_id,omitempty"` // ID from the user_groups table
}