package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"voicechat/internal/domain"
	"voicechat/internal/service"
)

type RoomHandler struct {
	roomService *service.RoomService
	seatService *service.SeatService
}

func NewRoomHandler(roomService *service.RoomService, seatService *service.SeatService) *RoomHandler {
	return &RoomHandler{roomService: roomService, seatService: seatService}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req service.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	room, err := h.roomService.CreateRoom(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	rooms, err := h.roomService.ListRooms(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	room, err := h.roomService.GetRoom(roomID)
	if err != nil {
		if err == domain.ErrRoomNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	userID := c.GetString("user_id")
	roomID := c.Param("room_id")

	err := h.roomService.DeleteRoom(roomID, userID)
	if err != nil {
		if err == domain.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not room owner"})
			return
		}
		if err == domain.ErrRoomNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}

func (h *RoomHandler) OccupySeat(c *gin.Context) {
	userID := c.GetString("user_id")
	roomID := c.Param("room_id")
	seatID := c.Param("seat_id")

	err := h.seatService.OccupySeat(seatID, userID, roomID)
	if err != nil {
		if err == domain.ErrSeatOccupied {
			c.JSON(http.StatusConflict, gin.H{"error": "seat is already occupied"})
			return
		}
		if err == domain.ErrSeatNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "seat not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "seat occupied"})
}

func (h *RoomHandler) VacateSeat(c *gin.Context) {
	seatID := c.Param("seat_id")

	err := h.seatService.VacateSeat(seatID)
	if err != nil {
		if err == domain.ErrSeatNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "seat not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "seat vacated"})
}
