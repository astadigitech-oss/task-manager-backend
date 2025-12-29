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

	pdfBytes, err := c.projectService.ExportProject(uint(projectID), currentUser.ID, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// --- START: NEW CONTROLLERS ADDED FOR NEW ROUTES ---

// Handler for Weekly Backward Report
func (c *ExportController) ExportWeeklyBackward(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	user, _ := ctx.Get("user")
	currentUser, _ := user.(*models.User)

	pdfBytes, err := c.projectService.ExportWeeklyBackward(uint(projectID), currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report_weekly_backward.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Handler for Weekly Forward Report
func (c *ExportController) ExportWeeklyForward(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	user, _ := ctx.Get("user")
	currentUser, _ := user.(*models.User)

	pdfBytes, err := c.projectService.ExportWeeklyForward(uint(projectID), currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report_weekly_forward.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Handler for Daily Report
func (c *ExportController) ExportDaily(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	user, _ := ctx.Get("user")
	currentUser, _ := user.(*models.User)

	pdfBytes, err := c.projectService.ExportDaily(uint(projectID), currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report_daily.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (c *ExportController) ExportAgenda(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	user, _ := ctx.Get("user")
	currentUser, _ := user.(*models.User)

	pdfBytes, err := c.projectService.ExportAgenda(uint(projectID), currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report_agenda.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
