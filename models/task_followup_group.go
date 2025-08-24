package models

// TaskFollowupGroup represents the follow-up of a task by a group
type TaskFollowupGroup struct {
	ID      uint  `gorm:"primaryKey"`
	GroupID uint  `gorm:"not null"` // FK to groups.id
	TaskID  uint  `gorm:"not null"` // FK to tasks.id
	Group   Group `gorm:"foreignKey:GroupID"` // Belongs to Group
}
