package routes

import (
	"project-management-backend/config"
	"project-management-backend/controllers"
	"project-management-backend/middleware"
	"project-management-backend/repositories"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	//Inisiasi Layer
	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	//repositories
	attendanceRepo := repositories.NewAttendanceRepository(config.DB)
	attendanceImageRepo := repositories.NewAttendanceImageRepository(config.DB)
	taskImageRepo := repositories.NewTaskImageRepository()
	taskFileRepo := repositories.NewTaskFileRepository(config.DB)
	taskRepo := repositories.NewTaskRepository()
	taskStatusLog := repositories.NewTaskStatusLogRepository()
	projectRepo := repositories.NewProjectRepository()
	projectImageRepo := repositories.NewProjectImageRepository()
	workspaceRepo := repositories.NewWorkspaceRepository()

	// Initialize Activity Logger
	activityLogger := utils.NewActivityLogger(config.DB)

	//services
	pdfService := services.NewPDFService()
	attendanceImageService := services.NewAttendanceImageService(attendanceImageRepo)
	attendanceService := services.NewAttendanceService(*attendanceRepo, *attendanceImageRepo, userRepo, workspaceRepo)
	projectService := services.NewProjectService(projectRepo, userRepo, workspaceRepo, taskRepo, taskStatusLog, pdfService, activityLogger) // Tambahkan userRepo dan pdfService
	taskImageService := services.NewTaskImageService(taskImageRepo, taskRepo, projectRepo, workspaceRepo)
	taskFileService := services.NewTaskFileService(taskFileRepo, taskRepo, projectRepo)
	taskService := services.NewTaskService(taskRepo, taskStatusLog, activityLogger)
	projectImageService := services.NewProjectImageService(projectImageRepo, projectRepo, workspaceRepo, userRepo)
	workspaceService := services.NewWorkspaceService(workspaceRepo, projectRepo, taskRepo)
	userService := services.NewUserService(userRepo)
	webSocketService := services.NewWebSocketService(userRepo, workspaceRepo, projectRepo, taskRepo)
	dashboardService := services.NewDashboardService(taskRepo)
	profileService := services.NewProfileService(userRepo)

	//controllers
	attendanceController := controllers.NewAttendanceController(*attendanceService, *attendanceImageService, pdfService, workspaceService)
	taskImageController := controllers.NewTaskImageController(taskImageService)
	taskFileController := controllers.NewTaskFileController(taskFileService)
	taskController := controllers.NewTaskController(taskService)
	projectController := controllers.NewProjectController(projectService)
	projectImageController := controllers.NewProjectImageController(projectImageService)
	workspaceController := controllers.NewWorkspaceController(workspaceService)
	userController := controllers.NewUserController(userService)
	webSocketController := controllers.NewWebSocketController(authService, webSocketService, userService)
	dashboardController := controllers.NewDashboardController(dashboardService)
	profileController := controllers.NewProfileController(profileService)
	exportController := controllers.NewExportController(projectService)

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
		api.GET("/profile", profileController.ListProfile)
		api.PUT("/profile", profileController.UpdateProfile)
		api.DELETE("/profile/image", profileController.DeleteProfileImage)

		// Dashboard
		api.GET("/dashboard", dashboardController.GetUserDashboard)
		api.GET("/dashboard/admin", dashboardController.GetAdminDashboard)

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

				// Attendance
				attendances := workspace.Group("/attendances")
				{
					attendances.POST("", attendanceController.SubmitAttendance)
					attendances.GET("/export", adminMiddleware, attendanceController.ExportAttendances)
				}
			}
		}

		// user
		users := api.Group("/users")
		{
			users.GET("", userController.GetAllUsers)
			users.DELETE("/delete/:user_id", adminMiddleware, userController.DeleteUser)
		}

		// Project
		projects := api.Group("/projects")
		{
			projects.GET("", projectController.ListProjects)
			projects.POST("", adminMiddleware, projectController.CreateProject)
			exportGroup := projects.Group("/:project_id/export")
			{
				// New specific routes
				exportGroup.GET("/daily", adminMiddleware, exportController.ExportDaily)
				exportGroup.GET("/weekly-backward", adminMiddleware, exportController.ExportWeeklyBackward)
				exportGroup.GET("/weekly-forward", adminMiddleware, exportController.ExportWeeklyForward)
				exportGroup.GET("/monitoring", adminMiddleware, exportController.ExportMonitoring)
			}

			project := projects.Group("/:project_id")
			{
				project.GET("", projectController.DetailProject)
				project.PUT("", adminMiddleware, projectController.UpdateProject)
				project.DELETE("", adminMiddleware, projectController.SoftDeleteProject)
				project.DELETE("/permanent", adminMiddleware, projectController.DeleteProject)

				project.GET("/members", projectController.GetMembers)
				project.POST("/members", adminMiddleware, projectController.AddMember)
				project.DELETE("/members/:user_id", adminMiddleware, projectController.RemoveSingleMember)
				project.DELETE("/members", adminMiddleware, projectController.RemoveMember)

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

				// Task Files
				files := task.Group("/files")
				{
					files.GET("", adminMiddleware, taskFileController.ListFiles)
					files.POST("", adminMiddleware, taskFileController.UploadFile)
					files.GET("/:fileId/view", adminMiddleware, taskFileController.ViewFile)
					files.GET("/:fileId/download", adminMiddleware, taskFileController.DownloadFile)
					files.DELETE("/:fileId", adminMiddleware, taskFileController.DeleteFile)
				}
			}
		}
	}

	r.Static("/uploads", "./uploads")
	r.Static("/profile-images", "./uploads/profile_images")
}
