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

func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	pc.Service.UpdateProfile(c)
}

func (pc *ProfileController) DeleteProfileImage(c *gin.Context) {
	pc.Service.DeleteProfileImage(c)
}
