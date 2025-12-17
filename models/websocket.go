package models

import (
	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
// This struct is now in its own file under the 'models' package
// for better code organization and separation of concerns.
// It holds a reference to the WebSocket hub, the connection itself,
// a channel for outbound messages, and user/workspace identifiers.

type Client struct {
	Hub         *Hub
	Conn        *websocket.Conn
	Send        chan []byte
	UserID      uint
	WorkspaceID uint
}

// Hub manages all WebSocket connections and broadcasts messages.
// It keeps track of all active clients and handles their registration,
// unregistration, and message broadcasting.
// This struct is central to the WebSocket communication system.

type Hub struct {
	Clients    map[uint]map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}
