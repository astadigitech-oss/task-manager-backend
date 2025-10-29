package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func ListTaskImages(c *gin.Context) {
	taskID := c.Param("task_id")
	var images []models.TaskImage
	if err := config.DB.Where("task_id = ?", taskID).Find(&images).Error; err != nil {
		utils.Error(0, "LIST_TASK_IMAGES", "task_images", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "LIST_TASK_IMAGES", "task_images", 0, "List task images")
	c.JSON(http.StatusOK, images)
}

func UploadTaskImage(c *gin.Context) {
	var input models.TaskImage
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "UPLOAD_TASK_IMAGE", "task_images", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.UploadedBy, "UPLOAD_TASK_IMAGE", "task_images", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(input.UploadedBy, "UPLOAD_TASK_IMAGE", "task_images", input.ID, input.URL)
	c.JSON(http.StatusCreated, input)
}
