package controllers

import (
	"log"
	"net/http"
	"project-management-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketController handles incoming WebSocket requests.

type WebSocketController struct {
	AuthService      services.AuthService
	WebSocketService services.WebSocketService
	UserService      services.UserService
}

// NewWebSocketController creates a new WebSocketController.
func NewWebSocketController(authService services.AuthService, webSocketService services.WebSocketService, userService services.UserService) *WebSocketController {
	return &WebSocketController{AuthService: authService, WebSocketService: webSocketService, UserService: userService}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// ServeWs handles WebSocket requests.
func (wsc *WebSocketController) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	token := c.Query("token")
	workspaceIDStr := c.Query("workspace_id")

	user, err := wsc.AuthService.GetUserFromToken(token)
	if err != nil {
		msg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid or expired token")
		conn.WriteMessage(websocket.CloseMessage, msg)
		conn.Close()
		return
	}

	if user.Role != "admin" {
		workspaceID, err := strconv.ParseUint(workspaceIDStr, 10, 64)
		if err != nil {
			msg := websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "Invalid workspace ID format")
			conn.WriteMessage(websocket.CloseMessage, msg)
			conn.Close()
			return
		}

		isMember, err := wsc.UserService.IsUserMemberOfWorkspace(user.ID, uint(workspaceID))
		if err != nil || !isMember {
			msg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Access Denied: Not a workspace member")
			conn.WriteMessage(websocket.CloseMessage, msg)
			conn.Close()
			return
		}

		wsc.WebSocketService.RegisterAndServeClient(conn, user.ID, uint(workspaceID))
	} else {
		var workspaceID uint64 = 0
		if workspaceIDStr != "" {
			parsedWorkspaceID, err := strconv.ParseUint(workspaceIDStr, 10, 64)
			if err != nil {
				msg := websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "Invalid workspace ID format")
				conn.WriteMessage(websocket.CloseMessage, msg)
				conn.Close()
				return
			}
			workspaceID = parsedWorkspaceID
		}
		wsc.WebSocketService.RegisterAndServeClient(conn, user.ID, uint(workspaceID))
	}

}
