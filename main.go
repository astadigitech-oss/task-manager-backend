package main

import (
	"log"
	"os"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.ConnectDB()

	// Auto-migrate semua model
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.Project{},
		&models.ProjectImage{},
		&models.Task{},
		&models.TaskImage{},
		&models.WorkspaceUser{},
		&models.ProjectUser{},
		&models.TaskUser{},
		&models.ActivityLog{},
		&models.ErrorLog{},
	)
	if err != nil {
		log.Fatal("Migration error: ", err)
	}

	router := gin.Default()
	routes.SetupRoutes(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
