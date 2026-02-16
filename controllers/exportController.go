package controllers

import (
	"net/http"
	"strconv"

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

// Handler for Weekly Backward Report
func (c *ExportController) ExportWeeklyBackward(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "project ID tidak ditemukan"})
		return
	}

	currentUser := GetCurrentUser(ctx)

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
	projectID, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "project ID tidak ditemukan"})
		return
	}

	currentUser := GetCurrentUser(ctx)

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
	projectID, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "project ID tidak ditemukan"})
		return
	}

	currentUser := GetCurrentUser(ctx)

	pdfBytes, err := c.projectService.ExportDaily(uint(projectID), currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report_daily.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (c *ExportController) ExportMonitoring(ctx *gin.Context) {
	projectID, err := strconv.Atoi(ctx.Param("project_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "project ID tidak ditemukan"})
		return
	}

	currentUser := GetCurrentUser(ctx)

	pdfBytes, err := c.projectService.ExportMonitoring(uint(projectID), currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=project_report_weekly_monitoring.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
