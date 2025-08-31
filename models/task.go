package models

import "time"

// Task represents the task model
type Task struct {
	ID             uint              `gorm:"primaryKey" json:"ID"`
	Label          string            `gorm:"not null" json:"Label"`
	TaskTypeID     uint              `gorm:"not null" json:"TaskTypeID"`
	Priority       string            `gorm:"type:enum('Normal', 'Medium', 'High', 'Escalation');default:'Normal'" json:"Priority"`
	StartDate      time.Time         `gorm:"type:date" json:"StartDate"`
	DueDate        *time.Time        `gorm:"type:date" json:"DueDate"`
	Description    string            `gorm:"type:longtext" json:"Description"`
	Attachment     string            `gorm:"type:varchar(255);nullable" json:"Attachment"`
	Status         string            `gorm:"type:enum('Pending','In Progress','In Review','Completed');default:'Pending'" json:"Status"`
	CreatedBy      uint              `gorm:"not null" json:"CreatedBy"`
	Creator        User              `gorm:"foreignKey:CreatedBy" json:"Creator"`
	CreatedAt      time.Time         `gorm:"type:timestamp;autoCreateTime" json:"CreatedAt"`
	UpdatedAt      time.Time         `gorm:"type:timestamp;autoUpdateTime" json:"UpdatedAt"`
	AssignedUsers  []AssignTaskToUser `gorm:"foreignKey:TaskID" json:"AssignedUsers"`
	AssignedGroups []AssignTaskToGroup `gorm:"foreignKey:TaskID" json:"AssignedGroups"`
	FollowupUsers  []TaskFollowupUser `gorm:"foreignKey:TaskID" json:"FollowupUsers"`
	FollowupGroups []TaskFollowupGroup `gorm:"foreignKey:TaskID" json:"FollowupGroups"`
	Comments       []TaskCommentLog  `gorm:"foreignKey:TaskID" json:"Comments"`
}