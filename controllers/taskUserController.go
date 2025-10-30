package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetTaskMembers(c *gin.Context) {
	var members []models.TaskUser
	if err := config.DB.Preload("User").Preload("Task").Find(&members).Error; err != nil {
		utils.Error(0, "GET_TASK_MEMBERS", "task_users", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func AddMemberToTask(c *gin.Context) {
	var input models.TaskUser
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "ADD_MEMBER_TO_TASK", "task_users", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.UserID, "ADD_MEMBER_TO_TASK", "task_users", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.ActivityLog(input.UserID, "ADD_MEMBER_TO_TASK", "task_users", input.ID, nil, input)
	c.JSON(http.StatusCreated, input)
}
