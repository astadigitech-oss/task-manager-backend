// controllers/project_image_controller.go
package controllers

import (
	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

type ProjectImageController struct {
	Service services.ProjectImageService
}

func NewProjectImageController(service services.ProjectImageService) *ProjectImageController {
	return &ProjectImageController{Service: service}
}

func (pic *ProjectImageController) GetProjectImages(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	images, err := pic.Service.GetProjectImages(projectID, currentUser.ID)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	var imageList []gin.H
	for _, image := range images {
		imageList = append(imageList, gin.H{
			"id":          image.ID,
			"url":         image.URL,
			"project_id":  image.ProjectID,
			"uploaded_by": image.UploadedBy,
			"created_at":  image.CreatedAt,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project images berhasil diambil",
		Data:    imageList,
	})
}

func (pic *ProjectImageController) UploadProjectImage(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
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

	// Upload image
	projectImage, err := pic.Service.UploadProjectImage(projectID, file, currentUser.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Image berhasil diupload",
		Data: gin.H{
			"id":          projectImage.ID,
			"url":         projectImage.URL,
			"project_id":  projectImage.ProjectID,
			"uploaded_by": projectImage.UploadedBy,
			"created_at":  projectImage.CreatedAt,
		},
	})
}

func (pic *ProjectImageController) DeleteProjectImage(c *gin.Context) {
	imageID, err := ParseUintParam(c, "image_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := pic.Service.DeleteProjectImage(imageID, currentUser.ID); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Image berhasil dihapus",
		Data:    gin.H{"image_id": imageID},
	})
}
