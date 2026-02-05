package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	Send        chan []byte
	UserID      uint
	WorkspaceID uint
}

type Hub struct {
	Clients    map[uint]map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}
