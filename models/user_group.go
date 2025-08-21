package models

// UserGroup represents the many-to-many relationship between users and groups
type UserGroup struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint `gorm:"not null"` // FK to users.id
	GroupID uint `gorm:"not null"` // FK to groups.id
}
