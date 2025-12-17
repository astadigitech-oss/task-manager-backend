package controllers

import (
	"fmt"
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register - Register user baru
func (ac *AuthController) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
		Position string `json:"position"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json", "auth", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	}

	if input.Position != "" {
		user.Position = &input.Position
	}

	if err := ac.AuthService.Register(user); err != nil {
		errorMsg := err.Error()
		utils.Error(0, "register", "auth", 0, errorMsg, "")

		if errorMsg == "email sudah terdaftar" {
			c.JSON(400, gin.H{
				"success": false,
				"code":    400,
				"error":   errorMsg,
			})
		} else if strings.Contains(errorMsg, "gagal memeriksa email") ||
			strings.Contains(errorMsg, "gagal mengenkripsi password") {
			fmt.Printf("Registration error: %v\n", err)
			c.JSON(500, gin.H{
				"success": false,
				"code":    500,
				"error":   "Terjadi kesalahan sistem",
			})
		} else {
			c.JSON(500, gin.H{
				"success": false,
				"code":    500,
				"error":   "Terjadi kesalahan sistem",
			})
		}
		return
	}

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "User berhasil didaftarkan",
		Data: gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"position": user.Position,
		},
	})
}

// Login - Login user
func (ac *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json", "auth", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, user, err := ac.AuthService.Login(input.Email, input.Password)
	if err != nil {
		utils.Error(0, "login", "auth", 0, err.Error(), "")
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"code":    200,
		"message": "Login berhasil",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"name":     user.Name,
				"email":    user.Email,
				"role":     user.Role,
				"position": user.Position,
			},
		},
	})
}

// GetProfile - Get user profile
func (ac *AuthController) GetProfile(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		utils.Error(0, "get_profile", "auth", 0, "User not authenticated", "")
		c.JSON(401, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	currentUser := user.(*models.User)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Profile berhasil diambil",
		Data: gin.H{
			"id":       currentUser.ID,
			"name":     currentUser.Name,
			"email":    currentUser.Email,
			"role":     currentUser.Role,
			"position": currentUser.Position,
		},
	})
}
