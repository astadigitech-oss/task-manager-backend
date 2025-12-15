package controllers

import (
	"errors"
	"log"
	"net/http"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	userID      uint
	workspaceID uint
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	clients    map[uint]map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.workspaceID] == nil {
				h.clients[client.workspaceID] = make(map[*Client]bool)
			}
			h.clients[client.workspaceID][client] = true
			updateUserOnlineStatus(client.userID, true)
			h.broadcastStatus(client.workspaceID, client.userID, "online")
		case client := <-h.unregister:
			if _, ok := h.clients[client.workspaceID][client]; ok {
				delete(h.clients[client.workspaceID], client)
				close(client.send)
				if len(h.clients[client.workspaceID]) == 0 {
					delete(h.clients, client.workspaceID)
				}
			}
			updateUserOnlineStatus(client.userID, false)
			h.broadcastStatus(client.workspaceID, client.userID, "offline")
		case message := <-h.broadcast:
			for _, clients := range h.clients {
				for client := range clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients[client.workspaceID], client)
					}
				}
			}
		}
	}
}

func (h *Hub) broadcastStatus(workspaceID, userID uint, status string) {
	message := []byte(`{"type": "user_status", "user_id": ` + strconv.Itoa(int(userID)) + `, "status": "` + status + `"}`)
	if clients, ok := h.clients[workspaceID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[workspaceID], client)
			}
		}
	}
}

var hub = newHub()

func ServeWs(c *gin.Context, authService services.AuthService) {
	go hub.run()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	token := c.Query("token")
	workspaceIDStr := c.Query("workspace_id")
	workspaceID, err := strconv.ParseUint(workspaceIDStr, 10, 64)
	if err != nil {
		log.Println("Invalid workspace ID format from client")
		msg := websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "Invalid workspace ID format")
		conn.WriteMessage(websocket.CloseMessage, msg)
		conn.Close()
		return
	}

	user, err := authService.GetUserFromToken(token)
	if err != nil {
		log.Println("Invalid token for websocket connection")
		msg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid or expired token")
		conn.WriteMessage(websocket.CloseMessage, msg)
		conn.Close()
		return
	}

	// Validate user is a member of the workspace
	var membership models.WorkspaceUser
	err = config.DB.Where("user_id = ? AND workspace_id = ?", user.ID, workspaceID).First(&membership).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Access denied: User %d is not a member of workspace %d", user.ID, workspaceID)
			msg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Access Denied: Not a workspace member")
			conn.WriteMessage(websocket.CloseMessage, msg)
			conn.Close()
			return
		}
		log.Printf("Database error checking membership: %v", err)
		msg := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Database error")
		conn.WriteMessage(websocket.CloseMessage, msg)
		conn.Close()
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), userID: user.ID, workspaceID: uint(workspaceID)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func updateUserOnlineStatus(userID uint, isOnline bool) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		log.Printf("Failed to find user %d: %v", userID, err)
		return
	}

	user.IsOnline = isOnline
	now := time.Now()
	if !isOnline {
		user.LastSeen = &now
	}

	if err := config.DB.Save(&user).Error; err != nil {
		log.Printf("Failed to update user %d status: %v", userID, err)
	}
}

// --- Custom Response Structs for Slimmer JSON ---
type SimpleWorkspaceResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SimpleProjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SimpleTaskResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type UserAssocResponse struct {
	ID         uint                      `json:"id"`
	Name       string                    `json:"name"`
	Email      string                    `json:"email"`
	Role       string                    `json:"role"`
	IsOnline   bool                      `json:"is_online"`
	LastSeen   *time.Time                `json:"last_seen,omitempty"`
	Workspaces []SimpleWorkspaceResponse `json:"workspaces"`
	Projects   []SimpleProjectResponse   `json:"projects"`
	Tasks      []SimpleTaskResponse      `json:"tasks"`
}

func GetOnlineUsers(c *gin.Context) {
	var users []models.User
	err := config.DB.
		Preload("Workspaces").
		Preload("Projects").
		Preload("Tasks.Task").
		Where("is_online = ?", true).
		Find(&users).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, []UserAssocResponse{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch online users"})
		return
	}

	var responseData []UserAssocResponse
	for _, user := range users {
		workspaces := []SimpleWorkspaceResponse{}
		for _, ws := range user.Workspaces {
			workspaces = append(workspaces, SimpleWorkspaceResponse{
				ID:   ws.ID,
				Name: ws.Name,
			})
		}

		projects := []SimpleProjectResponse{}
		for _, p := range user.Projects {
			projects = append(projects, SimpleProjectResponse{
				ID:   p.ID,
				Name: p.Name,
			})
		}

		tasks := []SimpleTaskResponse{}
		for _, taskUser := range user.Tasks {
			if taskUser.Task.ID != 0 {
				tasks = append(tasks, SimpleTaskResponse{
					ID:    taskUser.Task.ID,
					Title: taskUser.Task.Title,
				})
			}
		}

		userResponse := UserAssocResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
			IsOnline:   user.IsOnline,
			LastSeen:   user.LastSeen,
			Workspaces: workspaces,
			Projects:   projects,
			Tasks:      tasks,
		}

		responseData = append(responseData, userResponse)
	}

	c.JSON(http.StatusOK, responseData)
}
