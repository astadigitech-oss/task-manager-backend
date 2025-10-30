package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetProjectMembers(c *gin.Context) {
	var members []models.ProjectUser
	if err := config.DB.Preload("User").Preload("Project").Find(&members).Error; err != nil {
		utils.Error(0, "GET_PROJECT_MEMBERS", "project_users", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func AddMemberToProject(c *gin.Context) {
	var input models.ProjectUser
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "ADD_MEMBER_TO_PROJECT", "project_users", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.UserID, "ADD_MEMBER_TO_PROJECT", "project_users", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// utils.Activity(input.UserID, "ADD_MEMBER_TO_PROJECT", "project_users", input.ID, "User assigned to project")
	c.JSON(http.StatusCreated, input)
}
