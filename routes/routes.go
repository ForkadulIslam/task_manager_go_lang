package routes

import (
	"github.com/gin-gonic/gin"
	"taskmanager/controllers"
	"taskmanager/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/hello", controllers.Hello)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

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

		// Task routes
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.GetTasks)
		auth.GET("/tasks/:id", controllers.GetTaskByID)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)

		// UserGroup routes
		auth.POST("/user-groups", controllers.AddUserToGroup)
		auth.DELETE("/user-groups/:id", controllers.RemoveUserFromGroup)
		auth.GET("/groups/users/:group_id", controllers.GetUsersInGroup)
		auth.GET("/users/:user_id/groups", controllers.GetGroupsForUser)

		// AssignTaskToUser routes
		auth.POST("/assign-task-to-user", controllers.AssignTaskToUser)
		auth.DELETE("/assign-task-to-user/:id", controllers.RemoveTaskAssignmentFromUser)
		auth.GET("/users/:user_id/assigned-tasks", controllers.GetTasksAssignedToUser)
		auth.GET("/tasks/assigned-users/:task_id", controllers.GetUsersAssignedToTask)

		// AssignTaskToGroup routes
		auth.POST("/assign-task-to-group", controllers.AssignTaskToGroup)
		auth.DELETE("/assign-task-to-group/:id", controllers.RemoveTaskAssignmentFromGroup)
		auth.GET("/groups/assigned-tasks/:group_id", controllers.GetTasksAssignedToGroup)
		auth.GET("/tasks/assigned-groups/:task_id", controllers.GetGroupsAssignedToTask)

		// TaskFollowupUser routes
		auth.POST("/task-followup-users", controllers.AddTaskFollowupUser)
		auth.DELETE("/task-followup-users/:id", controllers.RemoveTaskFollowupUser)
		auth.GET("/tasks/followup-users/:task_id", controllers.GetFollowupUsersForTask)
	}

	r.POST("/test-date", controllers.TestDate)
