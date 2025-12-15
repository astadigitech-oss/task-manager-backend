package controllers

import (
	"fmt"
	"net/http"
	"project-management-backend/services"
	"project-management-backend/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{Service: service}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	search := c.DefaultQuery("search", "")
	role := c.DefaultQuery("role", "")
	position := c.DefaultQuery("position", "")

	// Convert to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// Prepare filters
	filters := services.UserFilters{
		Search:   strings.TrimSpace(search),
		Role:     strings.TrimSpace(role),
		Position: strings.TrimSpace(position),
		Page:     page,
		Limit:    limit,
	}

	result, err := uc.Service.GetAllUsersWithFilters(filters, currentUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Format response
	users := make([]gin.H, 0)
	for _, user := range result.Users {
		userData := gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"is_online":  user.IsOnline,
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// Add last_seen if available
		if user.LastSeen != nil {
			userData["last_seen"] = user.LastSeen.Format("2006-01-02 15:04:05")
		} else {
			userData["last_seen"] = nil
		}

		// Add position if exists
		if user.Position != nil && *user.Position != "" {
			userData["position"] = *user.Position
		}

		users = append(users, userData)
	}

	responseData := gin.H{
		"users": users,
		"pagination": gin.H{
			"total":       result.Total,
			"page":        result.Page,
			"limit":       result.Limit,
			"total_pages": result.TotalPages,
			"has_next":    result.Page < result.TotalPages,
			"has_prev":    result.Page > 1,
		},
	}

	// Add filters info if any filter is applied
	if search != "" || role != "" || position != "" {
		responseData["filters"] = gin.H{
			"search":   search,
			"role":     role,
			"position": position,
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("List user berhasil diambil (%d user)", len(users)),
		Data:    responseData,
	})
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	userID, err := ParseUintParam(c, "user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	user, err := uc.Service.GetUserByID(userID, currentUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	userData := gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.Position != nil && *user.Position != "" {
		userData["position"] = *user.Position
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: "Detail user berhasil diambil",
		Data:    userData,
	})
}

func (uc *UserController) SearchUsers(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query pencarian diperlukan"})
		return
	}

	currentUser := GetCurrentUser(c)

	users, err := uc.Service.SearchUsers(query, currentUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	userList := make([]gin.H, 0)
	for _, user := range users {
		userData := gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}

		if user.Position != nil && *user.Position != "" {
			userData["position"] = *user.Position
		}

		userList = append(userList, userData)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: "Hasil pencarian user",
		Data: gin.H{
			"query":       query,
			"users":       userList,
			"total_found": len(users),
		},
	})
}

func (uc *UserController) GetUserStats(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang bisa mengakses statistik"})
		return
	}

	// Get all users untuk statistik
	users, err := uc.Service.GetAllUsers(currentUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Hitung statistik
	stats := gin.H{
		"total_users": len(users),
		"roles":       make(map[string]int),
		"positions":   make(map[string]int),
	}

	for _, user := range users {
		// Count by role
		stats["roles"].(map[string]int)[user.Role]++

		// Count by position if exists
		if user.Position != nil && *user.Position != "" {
			position := *user.Position
			stats["positions"].(map[string]int)[position]++
		}
	}

	utils.ActivityLog(currentUser.ID, "VIEW_USER_STATS", "users", 0, nil, stats)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: "Statistik user berhasil diambil",
		Data:    stats,
	})
}

func (uc *UserController) GetUsersByIDs(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	if currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang bisa mengakses"})
		return
	}

	var input struct {
		UserIDs []uint `json:"user_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Limit jumlah user yang bisa di-request sekaligus
	if len(input.UserIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "maksimal 100 user per request"})
		return
	}

	users, err := uc.Service.GetUsersWithDetails(input.UserIDs, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userList := make([]gin.H, 0)
	for _, user := range users {
		userData := gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}

		if user.Position != nil && *user.Position != "" {
			userData["position"] = *user.Position
		}

		userList = append(userList, userData)
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Data %d user berhasil diambil", len(users)),
		Data: gin.H{
			"users":     userList,
			"requested": len(input.UserIDs),
			"found":     len(users),
			"not_found": len(input.UserIDs) - len(users),
		},
	})
}
