package main

import (
	"log"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
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
	router.Run(":8080")
}
