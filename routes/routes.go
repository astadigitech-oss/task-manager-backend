package routes

import (
	"project-management-backend/controllers"
	"project-management-backend/repositories"
	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	//Inisiasi Layer
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	workspaceRepo := repositories.NewWorkspaceRepository()
	workspaceService := services.NewWorkspaceService(workspaceRepo)
	workspaceController := controllers.NewWorkspaceController(workspaceService)

	projectRepo := repositories.NewProjectRepository()
	projectService := services.NewProjectService(projectRepo)
	projectController := controllers.NewProjectController(projectService)

	// User
	r.POST("/users", userController.CreateUser)
	r.GET("/users", userController.GetUsers)

	// Workspace
	workspaces := r.Group("/workspaces")
	{
		workspaces.GET("", workspaceController.ListWorkspaces)
		workspaces.POST("", workspaceController.CreateWorkspace)

		workspace := workspaces.Group("/:workspace_id")
		{
			workspace.GET("", workspaceController.DetailWorkspace)
			workspace.GET("/members", workspaceController.GetMembers)
			workspace.POST("/members", workspaceController.AddMember)
		}
	}

	// Project
	projects := r.Group("/workspaces/:workspace_id/projects")
	{
		projects.GET("", projectController.ListProjects)
		projects.POST("", projectController.CreateProject)

		project := projects.Group("/:project_id")
		{
			project.GET("", projectController.DetailProject)
			project.GET("/members", projectController.GetMembers)
			project.POST("/members", projectController.AddMember)
		}
	}

	// Task
	r.POST("/tasks", controllers.CreateTask)
	r.GET("/tasks", controllers.GetTasks)

	// ProjectImage
	r.POST("/project-images", controllers.UploadProjectImage)
	r.GET("/project-images/:project_id", controllers.ListProjectImages)

	// TaskImage
	r.POST("/task-images", controllers.UploadTaskImage)
	r.GET("/task-images/:task_id", controllers.ListTaskImages)

	// Activity Log

	// // Error Log

}
