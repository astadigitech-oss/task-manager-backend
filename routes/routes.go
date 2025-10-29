package routes

import (
	"project-management-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// User
	r.POST("/users", controllers.CreateUser)
	r.GET("/users", controllers.GetUsers)
	// r.GET("/users/:id", controllers.GetUserByID)
	// r.PUT("/users/:id", controllers.UpdateUser)
	// r.DELETE("/users/:id", controllers.DeleteUser)

	// Workspace
	r.POST("/workspaces", controllers.CreateWorkspace)
	r.GET("/workspaces", controllers.GetWorkspaces)
	// r.GET("/workspaces/:id", controllers.GetWorkspaceByID)
	// r.PUT("/workspaces/:id", controllers.UpdateWorkspace)
	// r.DELETE("/workspaces/:id", controllers.DeleteWorkspace)

	// Project
	r.POST("/projects", controllers.CreateProject)
	r.GET("/projects", controllers.GetProjects)
	// r.GET("/projects/:id", controllers.GetProjectByID)
	// r.PUT("/projects/:id", controllers.UpdateProject)
	// r.DELETE("/projects/:id", controllers.DeleteProject)

	// ProjectUser (assign member to project)
	r.POST("/project-members", controllers.AddMemberToProject)
	r.GET("/project-members", controllers.GetProjectMembers)
	// r.DELETE("/project-members/:id", controllers.RemoveMemberFromProject)

	// Task
	r.POST("/tasks", controllers.CreateTask)
	r.GET("/tasks", controllers.GetTasks)
	// r.GET("/tasks/:id", controllers.GetTaskByID)
	// r.PUT("/tasks/:id", controllers.UpdateTask)
	// r.DELETE("/tasks/:id", controllers.DeleteTask)

	// TaskUser (assign member to task)
	r.POST("/task-members", controllers.AddMemberToTask)
	r.GET("/task-members", controllers.GetTaskMembers)
	// r.DELETE("/task-members/:id", controllers.RemoveMemberFromTask)

	// ProjectImage
	r.POST("/project-images", controllers.UploadProjectImage)
	r.GET("/project-images/:project_id", controllers.ListProjectImages)

	// TaskImage
	r.POST("/task-images", controllers.UploadTaskImage)
	r.GET("/task-images/:task_id", controllers.ListTaskImages)

	// Activity Log
	// r.GET("/activity-logs", controllers.GetActivityLogs)

	// // Error Log
	// r.GET("/error-logs", controllers.GetErrorLogs)
}
