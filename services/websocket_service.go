package services

import (
	"project-management-backend/models"
	"project-management-backend/repositories"
	"project-management-backend/utils"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// WebSocketService handles WebSocket connections and communication.

type WebSocketService interface {
	RunHub()
	RegisterAndServeClient(conn *websocket.Conn, userID uint, workspaceID uint)
}

type webSocketService struct {
	hub           *models.Hub
	userRepo      repositories.UserRepository
	workspaceRepo repositories.WorkspaceRepository
	projectRepo   repositories.ProjectRepository
	taskRepo      repositories.TaskRepository
}

// NewWebSocketService creates a new WebSocketService.
func NewWebSocketService(userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository, projectRepo repositories.ProjectRepository, taskRepo repositories.TaskRepository) WebSocketService {
	hub := &models.Hub{
		Clients:    make(map[uint]map[*models.Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
	}
	return &webSocketService{
		hub:           hub,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
		projectRepo:   projectRepo,
		taskRepo:      taskRepo,
	}
}

// RunHub runs the WebSocket hub.
func (s *webSocketService) RunHub() {
	for {
		select {
		case client := <-s.hub.Register:
			if s.hub.Clients[client.WorkspaceID] == nil {
				s.hub.Clients[client.WorkspaceID] = make(map[*models.Client]bool)
			}
			s.hub.Clients[client.WorkspaceID][client] = true
			s.updateUserOnlineStatus(client.UserID, true)
			s.broadcastStatus(client.WorkspaceID, client.UserID, "online")

		case client := <-s.hub.Unregister:
			if _, ok := s.hub.Clients[client.WorkspaceID][client]; ok {
				delete(s.hub.Clients[client.WorkspaceID], client)
				close(client.Send)
				if len(s.hub.Clients[client.WorkspaceID]) == 0 {
					delete(s.hub.Clients, client.WorkspaceID)
				}
			}
			s.updateUserOnlineStatus(client.UserID, false)
			s.broadcastStatus(client.WorkspaceID, client.UserID, "offline")

		case message := <-s.hub.Broadcast:
			for _, clients := range s.hub.Clients {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						client.Hub.Unregister <- client
					}
				}
			}
		}
	}
}

// RegisterAndServeClient creates a client and starts serving it.
func (s *webSocketService) RegisterAndServeClient(conn *websocket.Conn, userID uint, workspaceID uint) {
	client := &models.Client{Hub: s.hub, Conn: conn, Send: make(chan []byte, 256), UserID: userID, WorkspaceID: workspaceID}
	s.hub.Register <- client

	go s.writePump(client)
	go s.readPump(client)
}

func (s *webSocketService) readPump(client *models.Client) {
	defer func() {
		client.Hub.Unregister <- client
		client.Conn.Close()
	}()
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Error(client.UserID, "read_message", "websocket", client.WorkspaceID, err.Error(), "")
			}
			break
		}
	}
}

func (s *webSocketService) writePump(client *models.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				utils.Error(client.UserID, "next_writer", "websocket", client.WorkspaceID, err.Error(), "")
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				utils.Error(client.UserID, "close_writer", "websocket", client.WorkspaceID, err.Error(), "")
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				utils.Error(client.UserID, "ping", "websocket", client.WorkspaceID, err.Error(), "")
				return
			}
		}
	}
}

func (s *webSocketService) broadcastStatus(workspaceID, userID uint, status string) {
	message := []byte(`{"type": "user_status", "user_id": ` + strconv.Itoa(int(userID)) + `, "status": "` + status + `"}`)
	if clients, ok := s.hub.Clients[workspaceID]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				client.Hub.Unregister <- client
			}
		}
	}
}

func (s *webSocketService) updateUserOnlineStatus(userID uint, isOnline bool) {
	if err := s.userRepo.UpdateUserOnlineStatus(userID, isOnline); err != nil {
		utils.Error(userID, "update_user_status", "users", userID, err.Error(), "")
	}
}
