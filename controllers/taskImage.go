// controllers/task_image_controller.go
package controllers

import (
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

type TaskImageController struct {
	Service services.TaskImageService
}

func NewTaskImageController(service services.TaskImageService) *TaskImageController {
	return &TaskImageController{Service: service}
}

func (tic *TaskImageController) GetTaskImages(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	images, err := tic.Service.GetTaskImages(taskID, projectID, workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	var imageList []gin.H
	for _, image := range images {
		imageList = append(imageList, gin.H{
			"id":          image.ID,
			"url":         image.URL,
			"task_id":     image.TaskID,
			"uploaded_by": image.UploadedBy,
			"created_at":  image.CreatedAt,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task images berhasil diambil",
		Data:    imageList,
	})
}

func (tic *TaskImageController) UploadTaskImage(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get file dari form-data
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "File image diperlukan"})
		return
	}

	currentUser := GetCurrentUser(c)

	taskImage, err := tic.Service.UploadTaskImage(taskID, workspaceID, file, currentUser)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "CREATE_TASK_IMAGE", "task", taskID, nil, taskImage)

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Image berhasil diupload",
		Data: gin.H{
			"id":          taskImage.ID,
			"url":         taskImage.URL,
			"task_id":     taskImage.TaskID,
			"uploaded_by": taskImage.UploadedBy,
			"created_at":  taskImage.CreatedAt,
		},
	})
}

// ✅ DELETE TASK IMAGE - Pattern sama dengan Project Image
func (tic *TaskImageController) DeleteTaskImage(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	imageID, err := ParseUintParam(c, "image_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := tic.Service.DeleteTaskImage(imageID, workspaceID, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "DELETE_TASK_IMAGE", "task_image", imageID, nil, "deleted image")

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Image berhasil dihapus",
		Data:    gin.H{"image_id": imageID},
	})
}
