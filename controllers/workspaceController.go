package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetWorkspaces(c *gin.Context) {
	var workspaces []models.Workspace
	if err := config.DB.Preload("Creator").Preload("Projects").Find(&workspaces).Error; err != nil {
		utils.Error(0, "GET_WORKSPACES", "workspaces", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "GET_WORKSPACES", "workspaces", 0, "Get all workspaces")
	c.JSON(http.StatusOK, workspaces)
}

func CreateWorkspace(c *gin.Context) {
	var input models.Workspace
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "CREATE_WORKSPACE", "workspaces", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.CreatedBy, "CREATE_WORKSPACE", "workspaces", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(input.CreatedBy, "CREATE_WORKSPACE", "workspaces", input.ID, input.Name)
	c.JSON(http.StatusCreated, input)
}
