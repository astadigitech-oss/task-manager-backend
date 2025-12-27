package controllers

import (
	"net/http"
	"strconv"

	"project-management-backend/models"
	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

type ExportController struct {
	projectService services.ProjectService
}

func NewExportController(projectService services.ProjectService) *ExportController {
	return &ExportController{
		projectService: projectService,
	}
}

func (c *ExportController) ExportProject(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	filter := ctx.Query("filter")

	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser, ok := user.(*models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	// Memperbaiki urutan argumen agar sesuai dengan service
	// yang diharapkan: ExportProject(projectID uint, userID uint, filter string)
	pdfBytes, err := c.projectService.ExportProject(uint(projectID), currentUser.ID, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
