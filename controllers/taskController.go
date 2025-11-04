// controllers/task_controller.go
package controllers

import (
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Service services.TaskService
}

func NewTaskController(service services.TaskService) *TaskController {
	return &TaskController{Service: service}
}

func (tc *TaskController) ListTasks(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)
	tasks, err := tc.Service.GetAllTasks(projectID, workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	respTasks := utils.ToTaskResponseList(tasks)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Task list diambil",
		Data:    respTasks,
	})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
		StartDate   time.Time `json:"start_date"`
		DueDate     time.Time `json:"due_date"`
		Notes       *string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	task := models.Task{
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		StartDate:   input.StartDate,
		DueDate:     input.DueDate,
		Notes:       input.Notes,
	}

	if err := tc.Service.CreateTask(&task, workspaceID, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "CREATE_TASK", "task", task.ID, nil, task)

	respTask := utils.ToTaskResponse(&task)
	c.JSON(201, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Task berhasil dibuat",
		Data:    respTask,
	})
}

func (tc *TaskController) AddMember(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role_in_task"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := tc.Service.AddMember(taskID, workspaceID, input.UserID, input.Role, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "ADD_MEMBER_TASK", "task", taskID, nil, input)

	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil ditambahkan ke task",
		Data: gin.H{
			"task_id": taskID,
			"user_id": input.UserID,
			"role":    input.Role,
		},
	})
}

func (tc *TaskController) GetMembers(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	members, err := tc.Service.GetMembers(taskID, workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	memberResponses := utils.ToTaskMemberResponseList(members)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Members task berhasil diambil",
		Data:    memberResponses,
	})
}

func (tc *TaskController) DetailTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	task, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	respTask := utils.ToTaskResponse(task)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Detail task berhasil diambil",
		Data:    respTask,
	})
}
