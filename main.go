package main

import (
	"log"
	"os"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CleanupOnlineStatus() {
	log.Println("Cleaning up user online status...")

	now := time.Now()
	result := config.DB.Model(&models.User{}).Where("is_online = ?", true).Updates(map[string]interface{}{
		"is_online": false,
		"last_seen": &now,
	})

	if result.Error != nil {
		log.Printf("Error cleaning up online status: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("Reset online status for %d user(s).", result.RowsAffected)
	} else {
		log.Println("No lingering online users found. Clean.")
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using system environment variables")
	}

	config.ConnectDB()

	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance")
	}
	defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(5)

	sqlDB.SetMaxOpenConns(20)

	sqlDB.SetConnMaxLifetime(30 * time.Hour)

	// Auto-migrate semua model
	err = config.DB.AutoMigrate(
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

	// Clean up any stale online statuses from a previous crash/restart
	CleanupOnlineStatus()

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
