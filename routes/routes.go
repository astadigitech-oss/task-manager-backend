package routes

import (
	"project-management-backend/controllers"
	"project-management-backend/middleware"
	"project-management-backend/repositories"
	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	//Inisiasi Layer
	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	//repositories
	taskRepo := repositories.NewTaskRepository()
	projectRepo := repositories.NewProjectRepository()
	projectImageRepo := repositories.NewProjectImageRepository()
	workspaceRepo := repositories.NewWorkspaceRepository()

	//services
	taskService := services.NewTaskService(taskRepo)
	projectService := services.NewProjectService(projectRepo, workspaceRepo)
	projectImageService := services.NewProjectImageService(projectImageRepo, projectRepo, workspaceRepo)
	workspaceService := services.NewWorkspaceService(workspaceRepo, projectRepo, taskRepo)

	//controllers
	taskController := controllers.NewTaskController(taskService)
	projectController := controllers.NewProjectController(projectService)
	projectImageController := controllers.NewProjectImageController(projectImageService)
	workspaceController := controllers.NewWorkspaceController(workspaceService)

	authMiddleware := middleware.AuthMiddleware(authService)
	adminMiddleware := middleware.AdminMiddleware()

	//public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/profile", authMiddleware, authController.GetProfile) // Ini butuh auth
	}

	//Protected Routes
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		// Workspace
		workspaces := api.Group("/workspaces")
		{
			workspaces.GET("", workspaceController.ListWorkspaces)
			workspaces.POST("", adminMiddleware, workspaceController.CreateWorkspace)

			workspace := workspaces.Group("/:workspace_id")
			{
				workspace.GET("", workspaceController.DetailWorkspace)
				workspace.PUT("", adminMiddleware, workspaceController.UpdateWorkspace)
				workspace.DELETE("", adminMiddleware, workspaceController.SoftDeleteWorkspace)
				workspace.DELETE("/permanent", adminMiddleware, workspaceController.DeleteWorkspace)

				workspace.GET("/members", adminMiddleware, workspaceController.GetMembers)
				workspace.POST("/members", adminMiddleware, workspaceController.AddMember)
			}
		}

		// Project
		projects := api.Group("projects")
		{
			projects.GET("", projectController.ListProjects)
			projects.POST("", adminMiddleware, projectController.CreateProject)

			project := projects.Group("/:project_id")
			{
				project.GET("", projectController.DetailProject)
				project.PUT("", adminMiddleware, projectController.UpdateProject)
				project.DELETE("", adminMiddleware, projectController.SoftDeleteProject)
				project.DELETE("/permanent", adminMiddleware, projectController.DeleteProject)

				project.GET("/members", adminMiddleware, projectController.GetMembers)
				project.POST("/members", adminMiddleware, projectController.AddMember)

				// Project Images
				images := project.Group("/images")
				{
					images.GET("", projectImageController.GetProjectImages)
					images.POST("", adminMiddleware, projectImageController.UploadProjectImage)
					images.DELETE("/:image_id", adminMiddleware, projectImageController.DeleteProjectImage)
				}
			}
		}

		// Task
		tasks := api.Group("/workspaces/:workspace_id/projects/:project_id/tasks")
		{
			tasks.GET("", taskController.ListTasks)
			tasks.POST("", adminMiddleware, taskController.CreateTask) // Member project

			task := tasks.Group("/:task_id")
			{
				task.GET("", taskController.DetailTask)
				task.PUT("", taskController.UpdateTask)                               // Member project/task
				task.DELETE("", adminMiddleware, taskController.SoftDeleteTask)       // Creator atau member project
				task.DELETE("/permanent", adminMiddleware, taskController.DeleteTask) // Creator atau member project

				task.GET("/members", adminMiddleware, taskController.GetMembers)
				task.POST("/members", adminMiddleware, taskController.AddMember) // Member project bisa add member task
			}
		}
	}

	// ==================== STATIC FILES ====================
	// Tetap public untuk akses uploaded images
	r.Static("/uploads", "./uploads")
}
