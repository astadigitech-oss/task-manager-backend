package controllers

import (
	"net/http"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	Service *services.DashboardService
}

func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{Service: service}
}

func (dc *DashboardController) GetUserDashboard(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	tasks, err := dc.Service.GetAllTasksForUser(currentUser.ID)
	if err != nil {
		utils.Error(currentUser.ID, "get_user_dashboard", "dashboard", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: "Dashboard data berhasil diambil",
		Data: gin.H{
			"tasks": utils.ToTaskResponseList(tasks),
		},
	})
}

func (dc *DashboardController) GetAdminDashboard(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	tasks, err := dc.Service.GetAllTaskForAdmin()
	if err != nil {
		utils.Error(currentUser.ID, "get_admin_dashboard", "dashboard", 0, err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard data"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    200,
		Message: "Dashboard data berhasil diambil",
		Data: gin.H{
			"tasks": utils.ToTaskResponseList(tasks),
		},
	})
}
