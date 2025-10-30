package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)
	roleVal, _ := c.Get("role")
	role := roleVal.(string)

	workspaceID := c.Query("workspaceId")

	var projects []models.Project
	query := config.DB.Preload("Workspace").
		Preload("Members.User").
		Preload("Tasks").
		Preload("Images")

	// Cuma ambil project yang diikuti/member (non-admin)
	if role != "admin" {
		query = query.Joins("JOIN project_users pu ON pu.project_id = projects.id").
			Where("pu.user_id = ?", userID)
	}

	// Filter workspace kalo dikirim
	if workspaceID != "" {
		query = query.Where("workspace_id = ?", workspaceID)
	}

	if err := query.Find(&projects).Error; err != nil {
		utils.Error(userID, "GET_PROJECTS", "projects", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Activity(userID, "GET_PROJECTS", "projects", 0, "Get all projects")
	c.JSON(http.StatusOK, projects)
}

func CreateProject(c *gin.Context) {
	userIDVal, userOk := c.Get("user_id")
	roleVal, roleOk := c.Get("role")
	userID, idOk := userIDVal.(uint)
	role, roleStr := roleVal.(string)

	// Validasi context
	if !userOk || !roleOk {
		utils.Error(userID, "GET_PROJECTS", "projects", 401, "Unauthorized", "")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if !idOk || !roleStr {
		utils.Error(userID, "GET_PROJECTS", "projects", 400, "Invalid user context", "")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user context"})
		return
	}

	// Hanya admin yang boleh create
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang boleh membuat project"})
		return
	}

	var input struct {
		WorkspaceID uint   `json:"workspaceId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Members     []uint `json:"members"` // ID user anggota
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(userID, "CREATE_PROJECT", "projects", 0, err.Error(), "Invalid JSON payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Buat project baru
	project := models.Project{
		WorkspaceID: input.WorkspaceID,
		Name:        input.Title,
		Description: input.Description,
		CreatedBy:   userID,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		utils.Error(userID, "CREATE_PROJECT", "projects", 0, err.Error(), "Database insert failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tambahkan members ke pivot ProjectUser (boleh kosong)
	for _, memberID := range input.Members {
		projectUser := models.ProjectUser{
			ProjectID:     project.ID,
			UserID:        memberID,
			RoleInProject: "", // Jabatan default null
		}
		config.DB.Create(&projectUser)
	}

	// Log aktivitas
	utils.Activity(userID, "CREATE_PROJECT", "projects", project.ID, project.Name)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Project berhasil dibuat",
		"project": project,
	})
}
