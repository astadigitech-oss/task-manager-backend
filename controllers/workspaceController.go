package controllers

import (
	"net/http"
	"project-management-backend/services"

	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	Service services.WorkspaceService
}

func NewWorkspaceController(service services.WorkspaceService) *WorkspaceController {
	return &WorkspaceController{Service: service}
}

// GET /api/workspaces
func (wc *WorkspaceController) GetWorkspaces(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint)

	workspaces, err := wc.Service.GetUserWorkspaces(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}

// POST /api/workspaces
func (wc *WorkspaceController) CreateWorkspace(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	roleVal, _ := c.Get("role")
	userID := userIDVal.(uint)
	role := roleVal.(string)

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	workspace, err := wc.Service.CreateWorkspace(userID, role, input.Name)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workspace)
}
