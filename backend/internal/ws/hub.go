package ws

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	clients    map[string]*Client
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	redis      *redis.Client
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	RoomID  string
	Event   string
	Payload interface{}
}

type Client struct {
	ID       string
	RoomID   string
	conn     interface{} // *websocket.Conn
	send     chan []byte
	hub      *Hub
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan *BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      redisClient,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("Client registered: %s\n", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
			h.mu.Unlock()
			fmt.Printf("Client unregistered: %s\n", client.ID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				if client.RoomID == msg.RoomID {
					payload, _ := json.Marshal(msg.Payload)
					message := Message{
						Event:   msg.Event,
						Payload: payload,
					}
					data, _ := json.Marshal(message)
					select {
					case client.send <- data:
					default:
						go func(c *Client) {
							h.unregister <- c
						}(client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID, event string, payload interface{}) {
	h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Event:   event,
		Payload: payload,
	}
}

func (h *Hub) SendToClient(clientID string, data []byte) bool {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	select {
	case client.send <- data:
		return true
	default:
		return false
	}
}

func (h *Hub) GetClientsInRoom(roomID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*Client
	for _, client := range h.clients {
		if client.RoomID == roomID {
			clients = append(clients, client)
		}
	}
	return clients
}
