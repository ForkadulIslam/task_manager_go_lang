package models

import "time"

// Task represents the task model
type Task struct {
	ID          uint      `gorm:"primaryKey"`
	Label       string    `gorm:"not null"`
	TaskTypeID  uint      `gorm:"not null"` // FK to task_types.id
	Priority    string    `gorm:"type:enum('Normal', 'Medium', 'High', 'Escalation');default:'Normal'"`
	StartDate   time.Time `gorm:"type:date"`
	DueDate     *time.Time `gorm:"type:date"`
	Description string    `gorm:"type:longtext"`
	Attachment  string    `gorm:"type:varchar(255);nullable"`
	Status      string    `gorm:"type:enum('Pending','In Progress','In Review','Completed');default:'Pending';comment:0=Pending,1=In Progress,2=In Review,3=Completed"`
	CreatedBy   uint      `gorm:"not null"` // FK to users.id
	Creator     User      `gorm:"foreignKey:CreatedBy"` // Add this line for preloading creator details
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
	AssignedUsers []AssignTaskToUser `gorm:"foreignKey:TaskID"` // One-to-Many relationship
	AssignedGroups []AssignTaskToGroup `gorm:"foreignKey:TaskID"` // One-to-Many relationship
	FollowupUsers []TaskFollowupUser `gorm:"foreignKey:TaskID"` // One-to-Many relationship
	FollowupGroups []TaskFollowupGroup `gorm:"foreignKey:TaskID" json:"FollowupGroups"` // One-to-Many relationship
	Comments      []TaskCommentLog   `gorm:"foreignKey:TaskID"`
}