package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project
	if err := config.DB.Preload("Workspace").Preload("Members.User").Preload("Tasks").Preload("Images").Find(&projects).Error; err != nil {
		utils.Error(0, "GET_PROJECTS", "projects", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "GET_PROJECTS", "projects", 0, "Get all projects")
	c.JSON(http.StatusOK, projects)
}

func CreateProject(c *gin.Context) {
	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "CREATE_PROJECT", "projects", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.CreatedBy, "CREATE_PROJECT", "projects", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(input.CreatedBy, "CREATE_PROJECT", "projects", input.ID, input.Name)
	c.JSON(http.StatusCreated, input)
}
