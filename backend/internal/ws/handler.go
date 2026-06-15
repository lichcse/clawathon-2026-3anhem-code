package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"voicechat/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for MVP
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, hub *Hub, userID string, db *sql.DB, memberRepo *repository.RoomMemberRepository, seatRepo *repository.SeatRepository) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	client := &Client{
		ID:   userID,
		conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
	}

	hub.register <- client

	// Send connected event
	connectedMsg := map[string]interface{}{
		"event": "connected",
		"payload": map[string]interface{}{
			"user_id":     userID,
			"server_time": time.Now().UnixMilli(),
		},
	}
	data, _ := json.Marshal(connectedMsg)
	client.send <- data

	// Handle client messages
	go func() {
		defer func() {
			conn.Close()
			hub.unregister <- client

			// Clean up: remove from room
			if client.RoomID != "" {
				memberRepo.Leave(client.RoomID, userID)
				seatRepo.VacateSeatsByUser(client.RoomID, userID)

				hub.BroadcastToRoom(client.RoomID, "user_left", map[string]interface{}{
					"user_id": userID,
					"left_at": time.Now().UnixMilli(),
				})
			}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket error: %v\n", err)
				}
				return
			}

			// Handle message
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			handleClientMessage(msg, client, hub, db, memberRepo, seatRepo)
		}
	}()

	// Send messages
	for message := range client.send {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func handleClientMessage(msg Message, client *Client, hub *Hub, db *sql.DB, memberRepo *repository.RoomMemberRepository, seatRepo *repository.SeatRepository) {
	var payload map[string]interface{}
	json.Unmarshal(msg.Payload, &payload)

	switch msg.Event {
	case "join_room":
		roomID, ok := payload["room_id"].(string)
		if !ok {
			return
		}

		client.RoomID = roomID
		memberRepo.Join(roomID, client.ID)

		// Send room snapshot
		seats, _ := seatRepo.GetSeatsForRoom(roomID)
		members, _ := memberRepo.GetMembers(roomID)

		seatDTOs := []map[string]interface{}{}
		for _, seat := range seats {
			seatDTOs = append(seatDTOs, map[string]interface{}{
				"id":            seat.ID,
				"row":           seat.RowNum,
				"col":           seat.ColNum,
				"occupied_by":   seat.OccupiedByID,
			})
		}

		memberDTOs := []map[string]interface{}{}
		for _, member := range members {
			memberDTOs = append(memberDTOs, map[string]interface{}{
				"user_id":  member.UserID,
				"is_muted": member.IsMuted,
			})
		}

		snapshot := map[string]interface{}{
			"event": "room_snapshot",
			"payload": map[string]interface{}{
				"room_id": roomID,
				"seats":   seatDTOs,
				"members": memberDTOs,
			},
		}
		data, _ := json.Marshal(snapshot)
		client.send <- data

		// Broadcast user joined
		hub.BroadcastToRoom(roomID, "user_joined", map[string]interface{}{
			"user_id":   client.ID,
			"joined_at": time.Now().UnixMilli(),
		})

	case "sit_down":
		seatID, ok := payload["seat_id"].(string)
		if !ok {
			return
		}

		if err := seatRepo.OccupySeat(seatID, client.ID); err == nil {
			hub.BroadcastToRoom(client.RoomID, "seat_occupied", map[string]interface{}{
				"seat_id": seatID,
				"user_id": client.ID,
			})
		}

	case "stand_up":
		seatID, ok := payload["seat_id"].(string)
		if !ok {
			return
		}

		if err := seatRepo.VacateSeat(seatID); err == nil {
			hub.BroadcastToRoom(client.RoomID, "seat_vacated", map[string]interface{}{
				"seat_id": seatID,
				"user_id": client.ID,
			})
		}

	case "mic_toggle":
		isMuted, ok := payload["is_muted"].(bool)
		if !ok {
			return
		}

		memberRepo.SetMuteStatus(client.RoomID, client.ID, isMuted)
		hub.BroadcastToRoom(client.RoomID, "user_mic_changed", map[string]interface{}{
			"user_id":   client.ID,
			"is_muted":  isMuted,
		})

	case "leave_room":
		if client.RoomID != "" {
			memberRepo.Leave(client.RoomID, client.ID)
			seatRepo.VacateSeatsByUser(client.RoomID, client.ID)

			hub.BroadcastToRoom(client.RoomID, "user_left", map[string]interface{}{
				"user_id": client.ID,
				"left_at": time.Now().UnixMilli(),
			})

			client.RoomID = ""
		}

	case "audio":
		// Forward audio to all users in room
		if client.RoomID != "" {
			audioData, ok := payload["data"].(string)
			if !ok {
				return
			}

			audioMsg := map[string]interface{}{
				"event": "audio",
				"payload": map[string]interface{}{
					"user_id": client.ID,
					"data":    audioData,
				},
			}
			data, _ := json.Marshal(audioMsg)

			for _, other := range hub.GetClientsInRoom(client.RoomID) {
				if other.ID != client.ID {
					select {
					case other.send <- data:
					default:
					}
				}
			}
		}
	}
}
