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
	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "List task berhasil di ambil",
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
	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Task berhasil dibuat",
		Data:    respTask,
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

func (tc *TaskController) UpdateTask(c *gin.Context) {
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

	oldTask, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	task := models.Task{
		ID:          taskID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		StartDate:   input.StartDate,
		DueDate:     input.DueDate,
		Notes:       input.Notes,
	}

	if err := tc.Service.UpdateTask(&task, workspaceID, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	// Get data baru setelah update
	updatedTask, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Gagal mengambil data task setelah update"})
		return
	}

	utils.ActivityLog(currentUser.ID, "UPDATE_TASK", "task", taskID, oldTask, updatedTask)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task berhasil diupdate",
		Data: gin.H{
			"id":          updatedTask.ID,
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"status":      updatedTask.Status,
			"priority":    updatedTask.Priority,
			"start_date":  updatedTask.StartDate,
			"due_date":    updatedTask.DueDate,
			"notes":       updatedTask.Notes,
			"project_id":  updatedTask.ProjectID,
			"created_at":  updatedTask.CreatedAt,
		},
	})
}

// SoftDeleteTask - Soft delete task
func (tc *TaskController) SoftDeleteTask(c *gin.Context) {
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

	if err := tc.Service.SoftDeleteTask(taskID, workspaceID, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task berhasil di soft delete",
		Data:    gin.H{"task_id": taskID},
	})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
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

	// Get task info untuk show warning
	task, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	var input struct {
		Confirm bool `json:"confirm"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || !input.Confirm {
		c.JSON(400, gin.H{
			"error":   "Konfirmasi diperlukan untuk hard delete",
			"warning": "Tindakan ini akan menghapus PERMANEN:",
			"data_akan_dihapus": gin.H{
				"task":          task.Title,
				"members_count": len(task.Members),
				"images_count":  len(task.Images),
			},
			"confirmation_required": true,
		})
		return
	}

	if err := tc.Service.DeleteTask(taskID, workspaceID, currentUser); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task dan semua data terkait berhasil dihapus permanen",
		Data: gin.H{
			"task_id":         taskID,
			"task_title":      task.Title,
			"deleted_members": len(task.Members),
			"deleted_images":  len(task.Images),
		},
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

	c.JSON(200, APIResponse{
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
	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Members task berhasil diambil",
		Data:    memberResponses,
	})
}
