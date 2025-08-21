package models

import "time"

// User represents the user model
// @Description User model for the application
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Status    int       `gorm:"default:1;comment:1=Active,0=Inactive"`
	UserLabel int       `gorm:"comment:1=Super Admin,2=User"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

type UserResponse struct {
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	Status        int    `json:"status"`
	UserLabel     int    `json:"user_label"`
	AssociationID uint   `json:"association_id,omitempty"` // ID from the user_groups table
}
