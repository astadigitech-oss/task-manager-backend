package main

import (
	"log"
	"time"

	"project-management-backend/config"
	"project-management-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	config.ConnectDB()

	err := SeedData(config.DB)
	if err != nil {
		log.Fatalf("Could not seed data: %v", err)
	}
	log.Println("Successfully seeded data")
}

func SeedData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Seed User
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		users := []models.User{
		    {
		        Name:     "Test User",
		        Email:    "test@example.com",
		        Password: string(hashedPassword),
		        Role:     "member",
		    },
		    {
		        Name:     "Admin",
		        Email:    "admin@example.com",
		        Password: string(hashedPassword),
		        Role:     "admin",
		    },
		}

if err := tx.Create(&users).Error; err != nil {
    return err
}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 2. Seed Workspace
		workspace := models.Workspace{
			Name:        "My First Workspace",
			Description: "This is a collaborative workspace for our team.",
			CreatedBy:   user.ID,
			Color:       "#4A90E2",
		}
		if err := tx.Create(&workspace).Error; err != nil {
			return err
		}

		// 3. Seed WorkspaceUser (linking user to workspace)
		workspaceUser := models.WorkspaceUser{
			WorkspaceID:     workspace.ID,
			UserID:          user.ID,
			RoleInWorkspace: &user.Role,
		}
		if err := tx.Create(&workspaceUser).Error; err != nil {
			return err
		}

		// 4. Seed Project
		project := models.Project{
			Name:        "Initial Project",
			Description: "This is the first project in our workspace.",
			WorkspaceID: workspace.ID,
			CreatedBy:   user.ID, // Link to the user who created it
		}
		if err := tx.Create(&project).Error; err != nil {
			return err
		}

		// 5. Seed ProjectUser (linking user to project)
		projectUser := models.ProjectUser{
			ProjectID:     project.ID,
			UserID:        user.ID,
			RoleInProject: "editor",
		}
		if err := tx.Create(&projectUser).Error; err != nil {
			return err
		}

		// 6. Seed Task
		task := models.Task{
			ProjectID:   project.ID,
			Title:       "Design the homepage",
			Description: "Create a mockup for the new homepage design.",
			Status:      "on_progress",
			Priority:    "high",
			StartDate:   time.Now(),
			DueDate:     time.Now().AddDate(0, 0, 7), // Due in 7 days
		}
		if err := tx.Create(&task).Error; err != nil {
			return err
		}

		// 7. Seed TaskUser (linking user to task)
		taskUser := models.TaskUser{
			TaskID:     task.ID,
			UserID:     user.ID,
			AssignedAt: time.Now(),
		}
		if err := tx.Create(&taskUser).Error; err != nil {
			return err
		}

		return nil // Commit the transaction
	})

	return err
}
