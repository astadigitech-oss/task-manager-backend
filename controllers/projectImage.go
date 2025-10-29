package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func ListProjectImages(c *gin.Context) {
	projectID := c.Param("project_id")
	var images []models.ProjectImage
	if err := config.DB.Where("project_id = ?", projectID).Find(&images).Error; err != nil {
		utils.Error(0, "LIST_PROJECT_IMAGES", "project_images", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "LIST_PROJECT_IMAGES", "project_images", 0, "List project images")
	c.JSON(http.StatusOK, images)
}

func UploadProjectImage(c *gin.Context) {
	var input models.ProjectImage
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "UPLOAD_PROJECT_IMAGE", "project_images", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.UploadedBy, "UPLOAD_PROJECT_IMAGE", "project_images", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(input.UploadedBy, "UPLOAD_PROJECT_IMAGE", "project_images", input.ID, input.URL)
	c.JSON(http.StatusCreated, input)
}
