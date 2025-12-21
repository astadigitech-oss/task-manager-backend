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
	taskImageRepo := repositories.NewTaskImageRepository()
	taskRepo := repositories.NewTaskRepository()
	projectRepo := repositories.NewProjectRepository()
	projectImageRepo := repositories.NewProjectImageRepository()
	workspaceRepo := repositories.NewWorkspaceRepository()

	//services
	taskImageService := services.NewTaskImageService(taskImageRepo, taskRepo, projectRepo, workspaceRepo)
	taskService := services.NewTaskService(taskRepo)
	projectService := services.NewProjectService(projectRepo, workspaceRepo)
	projectImageService := services.NewProjectImageService(projectImageRepo, projectRepo, workspaceRepo, userRepo)
	workspaceService := services.NewWorkspaceService(workspaceRepo, projectRepo, taskRepo)
	userService := services.NewUserService(userRepo)
	webSocketService := services.NewWebSocketService(userRepo, workspaceRepo, projectRepo, taskRepo)
	dashboardService := services.NewDashboardService(taskRepo)
	profileService := services.NewProfileService(userRepo)

	//controllers
	taskImageController := controllers.NewTaskImageController(taskImageService)
	taskController := controllers.NewTaskController(taskService)
	projectController := controllers.NewProjectController(projectService)
	projectImageController := controllers.NewProjectImageController(projectImageService)
	workspaceController := controllers.NewWorkspaceController(workspaceService)
	userController := controllers.NewUserController(userService)
	webSocketController := controllers.NewWebSocketController(authService, webSocketService, userService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	profileController := controllers.NewProfileController(profileService)

	authMiddleware := middleware.AuthMiddleware(authService)
	adminMiddleware := middleware.AdminMiddleware()

	// Jalankan WebSocket Hub
	go webSocketService.RunHub()

	//public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/profile", authMiddleware, authController.GetProfile) // Ini butuh auth
	}

	// web socket
	ws := r.Group("/ws")
	{
		ws.GET("", webSocketController.ServeWs)
	}

	//Protected Routes
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		//Profile
		api.PUT("/profile", profileController.UpdateProfile)
		api.DELETE("/profile/image", profileController.DeleteProfileImage)

		// Dashboard
		api.GET("/dashboard", dashboardController.GetUserDashboard)

		// Online Users
		api.GET("/online-users", adminMiddleware, userController.GetOnlineUsers)

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
				workspace.POST("/members", adminMiddleware, workspaceController.AddMembers)
				workspace.DELETE("/members/:user_id", adminMiddleware, workspaceController.RemoveSingleMember)
				workspace.DELETE("/members/", adminMiddleware, workspaceController.RemoveMember)

				workspace.GET("/online-members", userController.GetOnlineWorkspaceMembers)
			}
		}

		// user
		users := api.Group("/users")
		{
			users.GET("", userController.GetAllUsers)
		}

		// Project
		projects := api.Group("/projects")
		{
			projects.GET("", projectController.ListProjects)
			projects.POST("", adminMiddleware, projectController.CreateProject)

			project := projects.Group("/:project_id")
			{
				project.GET("", projectController.DetailProject)
				project.PUT("", adminMiddleware, projectController.UpdateProject)
				project.DELETE("", adminMiddleware, projectController.SoftDeleteProject)
				project.DELETE("/permanent", adminMiddleware, projectController.DeleteProject)

				project.GET("/members", projectController.GetMembers)
				project.POST("/members", adminMiddleware, projectController.AddMember)
				project.DELETE("/members/:user_id", adminMiddleware, projectController.RemoveSingleMember)
				project.DELETE("/members/", adminMiddleware, projectController.RemoveMember)

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
			tasks.POST("", adminMiddleware, taskController.CreateTask)

			task := tasks.Group("/:task_id")
			{
				task.GET("", taskController.DetailTask)
				task.PUT("", taskController.UpdateTask)
				task.DELETE("", adminMiddleware, taskController.SoftDeleteTask)
				task.DELETE("/permanent", adminMiddleware, taskController.DeleteTask)

				task.GET("/members", taskController.GetMembers)
				task.POST("/members", adminMiddleware, taskController.AddMember)
				task.DELETE("/members/:user_id", adminMiddleware, taskController.DeleteMember)

				// Task Images
				images := task.Group("/images")
				{
					images.GET("", taskImageController.GetTaskImages)
					images.POST("", taskImageController.UploadTaskImage)
					images.DELETE("/:image_id", taskImageController.DeleteTaskImage)
				}
			}
		}
	}

	r.Static("/uploads", "./uploads")
}
