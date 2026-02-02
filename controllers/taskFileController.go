package controllers

import (
	"net/http"
	"strconv"

	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

type TaskFileController struct {
	taskFileService *services.TaskFileService
}

func NewTaskFileController(taskFileService *services.TaskFileService) *TaskFileController {
	return &TaskFileController{taskFileService: taskFileService}
}

func (c *TaskFileController) UploadFile(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("task_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	currentUser := GetCurrentUser(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	taskFile, err := c.taskFileService.UploadFile(uint(taskID), currentUser.ID, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	ctx.JSON(http.StatusOK, taskFile)
}

func (c *TaskFileController) ListFiles(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("task_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	files, err := c.taskFileService.GetFilesByTaskID(uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	ctx.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "File berhasil diambil",
		Data:    files,
	})
}

func (c *TaskFileController) DownloadFile(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	fileData, mimeType, filename, err := c.taskFileService.DownloadFile(uint(fileID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	ctx.Header("Content-Type", mimeType)
	ctx.Header("Content-Length", strconv.Itoa(len(fileData)))

	ctx.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "File berhasil diunduh",
		Data: gin.H{
			"filename":  filename,
			"mime_type": mimeType,
		},
	})
}

func (c *TaskFileController) ViewFile(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	fileData, mimeType, _, err := c.taskFileService.DownloadFile(uint(fileID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.Header("Content-Disposition", "inline")
	ctx.Header("Content-Type", mimeType)
	ctx.Header("Content-Length", strconv.Itoa(len(fileData)))

	ctx.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "File berhasil ditampilkan",
		Data: gin.H{
			"filedata":  fileData,
			"mime_type": mimeType,
		},
	})
}

func (c *TaskFileController) DeleteFile(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID := uint(ctx.GetUint("user_id"))

	err = c.taskFileService.DeleteFile(uint(fileID), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	ctx.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "File berhasil dihapus",
		Data:    fileID,
	})
}
