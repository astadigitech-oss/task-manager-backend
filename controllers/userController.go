package controllers

import (
	"net/http"
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.service.GetAllUsers()
	if err != nil {
		utils.Error(0, "GET_USERS", "users", 500, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Password string  `json:"password"`
		Role     string  `json:"role"`
		Position *string `json:"position"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "CREATE_USER", "users", 400, err.Error(), "Invalid JSON payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
		Position: input.Position,
	}

	if err := uc.service.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.Activity(user.ID, "CREATE_USER", "users", user.ID, user.Name)
	c.JSON(http.StatusCreated, gin.H{"message": "Create User success"})
}
