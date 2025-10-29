package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := config.DB.Preload("Project").Preload("Members.User").Preload("Images").Find(&tasks).Error; err != nil {
		utils.Error(0, "GET_TASKS", "tasks", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "GET_TASKS", "tasks", 0, "Get all tasks")
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "CREATE_TASK", "tasks", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(0, "CREATE_TASK", "tasks", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "CREATE_TASK", "tasks", input.ID, input.Title)
	c.JSON(http.StatusCreated, input)
}
