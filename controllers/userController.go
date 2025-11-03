package controllers

import (
	"errors"
	"net/http"
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	respUsers := utils.ToUserResponseList(users)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "User list diambil",
		Data:    respUsers,
	})
	// c.JSON(http.StatusOK, utils.NewResponse(true, 200, "Data berhasil diambil", users))
}

func (uc *UserController) GetUsersByEmail(c *gin.Context) {
	email := c.Query("email")
	users, err := uc.service.GetUsersByEmail(email)
	if err != nil {
		code := 500
		msg := "Gagal mengambil user dengan email " + email
		// Jika error not found, jadikan 404
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = 404
			msg = "User dengan email " + email + " tidak ditemukan"
		}
		resp := utils.APIResponse{
			Success: false,
			Code:    code,
			Message: msg,
			Data:    nil,
		}
		c.JSON(code, resp)
		return
	}
	resp := utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "User dengan email " + email + " berhasil diambil",
		Data:    users,
	}
	c.JSON(200, resp)
}

func (uc *UserController) GetUsersByRole(c *gin.Context) {
	role := c.Query("role")
	users, err := uc.service.GetUsersByRole(role)
	if err != nil {
		resp := utils.APIResponse{
			Success: false,
			Code:    500,
			Message: "Gagal mengambil user dengan role " + role,
			Data:    nil,
		}
		c.JSON(500, resp)
		return
	}
	resp := utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "User dengan role " + role + " berhasil diambil",
		Data:    users,
	}
	c.JSON(200, resp)
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

	// CreateUser expects a value (models.User), so pass user (non-pointer).
	if err := uc.service.CreateUser(&user); err != nil {
		utils.Error(0, "CREATE_USER", "users", 400, err.Error(), "Failed to create user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Langsung log gunakan struct user, tanpa query ulang
	utils.ActivityLog(user.ID, "CREATE_USER", "users", user.ID, nil, user)

	respUser := utils.ToUserResponse(&user)
	c.JSON(201, utils.APIResponse{
		Success: true,
		Code:    201,
		Message: "User berhasil dibuat",
		Data:    respUser,
	})
}
