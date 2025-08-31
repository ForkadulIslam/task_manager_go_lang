package routes

import (
	"taskmanager/controllers"
	"taskmanager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/hello", controllers.Hello)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/sync-user", controllers.SyscUser)

	// Authenticated routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Group routes
		auth.POST("/groups", controllers.CreateGroup)
		auth.GET("/groups", controllers.GetGroups)
		auth.GET("/groups/:id", controllers.GetGroupByID)
		auth.PUT("/groups/:id", controllers.UpdateGroup)
		auth.DELETE("/groups/:id", controllers.DeleteGroup)

		// TaskType routes
		auth.POST("/task-types", controllers.CreateTaskType)
		auth.GET("/task-types", controllers.GetTaskTypes)
		auth.GET("/task-types/:id", controllers.GetTaskTypeByID)
		auth.PUT("/task-types/:id", controllers.UpdateTaskType)
		auth.DELETE("/task-types/:id", controllers.DeleteTaskType)

		// User routes
		auth.GET("/users", controllers.GetUsers)

		// Task routes
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.GetTasks)
		auth.GET("/my-tasks", controllers.GetMyTasks)
		auth.POST("/my-tasks/filter", controllers.GetMyTasksFiltered) // New route for filtered My Tasks
		auth.GET("/tasks/:id", controllers.GetTaskByID)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)
		auth.POST("/tasks/:id/status", controllers.UpdateTaskStatus)
		auth.POST("/tasks/:id/comments", controllers.AddTaskComment)

		// UserGroup routes
		auth.POST("/user-groups", controllers.AssignUsersToGroup)
		auth.GET("/user-groups", controllers.GetGroupsCreatedByUser)
		auth.DELETE("/user-groups/:id", controllers.RemoveUserFromGroup)
		auth.GET("/user-groups/my-groups", controllers.GetGroupsForUser)

		// Attachment upload route
		auth.POST("/upload-attachment", controllers.UploadAttachment)

	}

}
