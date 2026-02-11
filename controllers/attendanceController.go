package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	service              services.AttendanceService
	attendanceImgService services.AttendanceImageService
	pdfService           services.PDFService
	workspaceService     services.WorkspaceService
}

func NewAttendanceController(
	service services.AttendanceService,
	attendanceImgService services.AttendanceImageService,
	pdfService services.PDFService,
	workspaceService services.WorkspaceService,
) *AttendanceController {
	return &AttendanceController{
		service:              service,
		attendanceImgService: attendanceImgService,
		pdfService:           pdfService,
		workspaceService:     workspaceService,
	}
}

func (c *AttendanceController) SubmitAttendance(ctx *gin.Context) {
	workspaceIDStr := ctx.Param("workspace_id")
	workspaceID, err := strconv.ParseUint(workspaceIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	currentUser := GetCurrentUser(ctx)

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parsing form: %v", err)})
		return
	}

	attendance := models.Attendance{
		WorkspaceID: uint(workspaceID),
		UserID:      currentUser.ID,
		Activity:    form.Value["activity"][0],
		ClockIn:     time.Now(),
	}

	if len(form.Value["obstacle"]) > 0 {
		obstacle := form.Value["obstacle"][0]
		attendance.Obstacle = &obstacle
	}

	if err := c.service.SubmitAttendance(&attendance); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to submit attendance: %v", err)})
		return
	}

	files := form.File["images"]
	for _, file := range files {
		if _, err := c.attendanceImgService.UploadImage(attendance.ID, file); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload image: %v", err)})
			return
		}
	}

	createdAttendance, err := c.service.GetAttendanceByID(attendance.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created attendance"})
		return
	}
	response := utils.ToAttendanceResponse(*createdAttendance)

	ctx.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Attendance submitted successfully",
		Data:    response,
	})
}

func (c *AttendanceController) ExportAttendances(ctx *gin.Context) {
	workspaceIDStr := ctx.Param("workspace_id")
	workspaceID, err := strconv.ParseUint(workspaceIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	currentUser := GetCurrentUser(ctx)
	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date query parameter is required"})
		return
	}

	workspace, err := c.workspaceService.GetByID(uint(workspaceID), currentUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workspace details"})
		return
	}

	attendances, err := c.service.GetAttendancesForExport(uint(workspaceID), date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendances for export"})
		return
	}

	pdfBytes, err := c.pdfService.CreateAttendanceReportPDF(attendances, workspace.Name, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("attendance_report_%s_%s.pdf", workspace.Name, date)
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
