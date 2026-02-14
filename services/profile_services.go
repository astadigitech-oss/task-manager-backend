package services

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-management-backend/models"
	"project-management-backend/repositories"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProfileService struct {
	UserRepo repositories.UserRepository
}

func NewProfileService(userRepo repositories.UserRepository) *ProfileService {
	return &ProfileService{UserRepo: userRepo}
}

func (s *ProfileService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
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

	position := c.PostForm("position")
	if position != "" {
		u.Position = &position
	}

	file, err := c.FormFile("profile_image")
	if err == nil {
		oldProfileImage := u.ProfileImage

		uploadDir, _ := filepath.Abs("uploads/profile_images")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}

		filename := filepath.Base(file.Filename)
		serverPath := filepath.Join(uploadDir, filename)
		urlPath := fmt.Sprintf("/profile-images/%s", filename)

		if err := c.SaveUploadedFile(file, serverPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		u.ProfileImage = &urlPath

		if oldProfileImage != nil && *oldProfileImage != "" {
			go func(oldURLPath string) {
				oldServerPath := filepath.Join("uploads", strings.TrimPrefix(oldURLPath, "/"))
				if err := os.Remove(oldServerPath); err != nil {
					log.Printf("[Cleanup] Failed to delete %s: %s", oldServerPath, err.Error())
				} else {
					log.Printf("[Cleanup] Deleted %s", oldServerPath)
				}
			}(*oldProfileImage)
		}
	}

	if err := s.UserRepo.UpdateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, u)
}

func (s *ProfileService) DeleteProfileImage(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	u := user.(*models.User)

	if u.ProfileImage != nil && *u.ProfileImage != "" {
		serverPath := filepath.Join("uploads", strings.TrimPrefix(*u.ProfileImage, "/"))
		if err := os.Remove(serverPath); err != nil {
			log.Printf("[Delete] Failed to delete image file %s: %s", serverPath, err.Error())
		}
	}

	u.ProfileImage = nil
	if err := s.UserRepo.UpdateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile image deleted successfully"})
}
