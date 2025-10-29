package controllers

import (
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

// Get all users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Projects").Preload("Tasks").Find(&users).Error; err != nil {
		utils.Error(0, "GET_USERS", "users", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(0, "GET_USERS", "users", 0, "Get all users")
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "CREATE_USER", "users", 0, err.Error(), "")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		utils.Error(input.ID, "CREATE_USER", "users", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Activity(input.ID, "CREATE_USER", "users", input.ID, input.Name)
	c.JSON(http.StatusCreated, input)
}
