package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// ChatMessage represents a message sent between admin and user
type ChatMessage struct {
	Type      string `json:"type"`
	TicketID  string `json:"ticket_id"`
	SenderID  string `json:"sender_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	IsAdmin   bool   `json:"is_admin"`
	AdminName string `json:"admin_name,omitempty"`
}

// Client represents a WebSocket client
type Client struct {
	ID       string
	TicketID string
	IsAdmin  bool
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	mu       sync.Mutex
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[string]map[string]*Client // ticketID -> clientID -> Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
}

// AdminClient represents an admin WebSocket client for global notifications
type AdminClient struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	mu   sync.Mutex
}

// AdminHub maintains admin clients and broadcasts notifications
type AdminHub struct {
	clients    map[string]*AdminClient // clientID -> AdminClient
	Register   chan *AdminClient
	Unregister chan *AdminClient
	Broadcast  chan []byte
	mu         sync.RWMutex
}

// BroadcastMessage contains message and target room
type BroadcastMessage struct {
	TicketID string
	Message  []byte
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// NewAdminHub creates a new AdminHub
func NewAdminHub() *AdminHub {
	return &AdminHub{
		clients:    make(map[string]*AdminClient),
		Register:   make(chan *AdminClient),
		Unregister: make(chan *AdminClient),
		Broadcast:  make(chan []byte, 256),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.TicketID] == nil {
				h.clients[client.TicketID] = make(map[string]*Client)
			}
			h.clients[client.TicketID][client.ID] = client
			log.Printf("[Hub] Client registered: %s for ticket: %s", client.ID, client.TicketID)
			h.mu.Unlock()

			// Send connected message
			msg := ChatMessage{
				Type:     "connected",
				TicketID: client.TicketID,
				SenderID: client.ID,
				IsAdmin:  client.IsAdmin,
			}
			if data, err := json.Marshal(msg); err == nil {
				h.Broadcast <- &BroadcastMessage{TicketID: client.TicketID, Message: data}
			}

		case client := <-h.Unregister:
			h.mu.Lock()
			if room, ok := h.clients[client.TicketID]; ok {
				if _, exists := room[client.ID]; exists {
					delete(room, client.ID)
					close(client.Send)
					log.Printf("[Hub] Client unregistered: %s from ticket: %s", client.ID, client.TicketID)
					if len(room) == 0 {
						delete(h.clients, client.TicketID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			if room, ok := h.clients[msg.TicketID]; ok {
				for _, client := range room {
					select {
					case client.Send <- msg.Message:
					default:
						close(client.Send)
						delete(room, client.ID)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(ticketID string, message []byte) {
	h.Broadcast <- &BroadcastMessage{
		TicketID: ticketID,
		Message:  message,
	}
}

// WritePump sends messages from hub to WebSocket
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		c.mu.Lock()
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		c.mu.Unlock()

		if err != nil {
			log.Printf("[WritePump] Error writing to client %s: %v", c.ID, err)
			return
		}
	}
}

// Run starts the admin hub
func (ah *AdminHub) Run() {
	for {
		select {
		case client := <-ah.Register:
			ah.mu.Lock()
			ah.clients[client.ID] = client
			ah.mu.Unlock()
		case client := <-ah.Unregister:
			ah.mu.Lock()
			if _, ok := ah.clients[client.ID]; ok {
				delete(ah.clients, client.ID)
				close(client.Send)
			}
			ah.mu.Unlock()
		case msg := <-ah.Broadcast:
			ah.mu.RLock()
			for _, client := range ah.clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(ah.clients, client.ID)
				}
			}
			ah.mu.RUnlock()
		}
	}
}

// WritePump sends messages from admin hub to WebSocket
func (ac *AdminClient) WritePump() {
	defer ac.Conn.Close()
	for message := range ac.Send {
		ac.mu.Lock()
		err := ac.Conn.WriteMessage(websocket.TextMessage, message)
		ac.mu.Unlock()
		if err != nil {
			return
		}
	}
}
