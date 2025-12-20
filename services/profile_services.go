package services

import (
	"net/http"
	"os"
	"path/filepath"
	"project-management-backend/models"
	"project-management-backend/repositories"

	"github.com/gin-gonic/gin"
)

type ProfileService struct {
	UserRepo repositories.UserRepository
}

func NewProfileService(userRepo repositories.UserRepository) *ProfileService {
	return &ProfileService{UserRepo: userRepo}
}
func (s *ProfileService) UpdateProfile(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	u, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	name := c.PostForm("name")
	if name != "" {
		u.Name = name
	}

	file, err := c.FormFile("profile_image")
	if err == nil {
		filename := filepath.Base(file.Filename)
		path := "uploads/profiles/" + filename

		if err := os.MkdirAll("uploads/profiles", os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		u.ProfileImage = &path
	}

	if err := s.UserRepo.UpdateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, u)
}

func (s *ProfileService) DeleteProfileImage(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	u := user.(models.User)

	if u.ProfileImage != nil {
		if err := os.Remove(*u.ProfileImage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image file"})
			return
		}
	}

	u.ProfileImage = nil
	if err := s.UserRepo.UpdateUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile image deleted successfully"})
}
