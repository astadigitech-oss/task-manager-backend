package main

import (
	"log"
	"os"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using system environment variables")
	}

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

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           5 * 60,
	}))

	routes.SetupRoutes(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
