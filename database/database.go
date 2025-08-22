package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"taskmanager/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/task_manager?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database! \n" + err.Error())
	}

			database.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.TaskType{},
		&models.Task{},
		&models.UserGroup{},
		&models.AssignTaskToUser{},
		&models.AssignTaskToGroup{},
		&models.TaskFollowupUser{},
		&models.TaskStatusUpdateLog{},
		&models.TaskCommentLog{},
		&models.TaskSeenByUser{},
	)

	DB = database
}

