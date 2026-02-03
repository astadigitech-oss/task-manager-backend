package controllers

import (
	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	Service *services.ProfileService
}

func NewProfileController(service *services.ProfileService) *ProfileController {
	return &ProfileController{Service: service}
}

func (pc *ProfileController) ListProfile(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	profile, err := pc.Service.GetProfile(currentUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get profile data"})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Profile berhasil diambil",
		Data: gin.H{
			"name":     profile.Name,
			"email":    profile.Email,
			"position": profile.Position,
			"avatar":   profile.ProfileImage,
		},
	})
}

func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	pc.Service.UpdateProfile(c)
}

func (pc *ProfileController) DeleteProfileImage(c *gin.Context) {
	pc.Service.DeleteProfileImage(c)
}
