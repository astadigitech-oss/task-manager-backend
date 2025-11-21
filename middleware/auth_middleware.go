package middleware

import (
	"project-management-backend/models"
	"project-management-backend/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token dari header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header diperlukan"})
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Format token tidak valid. Gunakan: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token dan get user
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Token tidak valid: " + err.Error()})
			c.Abort()
			return
		}

		// Set user ke context
		c.Set("currentUser", user)
		c.Next()
	}
}

// AdminMiddleware - Hanya untuk admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if !exists {
			c.JSON(401, gin.H{"error": "User tidak terautentikasi"})
			c.Abort()
			return
		}

		currentUser := user.(*models.User)
		if currentUser.Role != "admin" {
			c.JSON(403, gin.H{"error": "Hanya admin yang boleh mengakses"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware - Untuk route yang bisa diakses dengan atau tanpa auth
func OptionalAuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				user, err := authService.GetUserFromToken(tokenString)
				if err == nil {
					c.Set("currentUser", user)
				}
			}
		}
		c.Next()
	}
}
