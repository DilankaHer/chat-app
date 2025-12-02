package handler

import (
	"duhchat/internal/repo"
	"duhchat/internal/ws"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type CreateRoom struct {
	RoomName string `json:"roomName" validate:"required"`
}

type RoomHandler struct {
	roomRepository    repo.RoomRepository
	messageRepository repo.MessageRepository
	hub               *ws.Hub
}

func NewRoomHandler(roomRepository repo.RoomRepository, messageRepository repo.MessageRepository, hub *ws.Hub) *RoomHandler {
	return &RoomHandler{roomRepository: roomRepository, messageRepository: messageRepository, hub: hub}
}

func (rr *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")

	claims := r.Context().Value("user").(jwt.MapClaims)
	if claims == nil {
		http.Error(w, "Can't find userId ", http.StatusInternalServerError)
		return
	}
	userId := claims["userId"].(string)
	username := claims["username"].(string)

	userRoom := &repo.UserRoom{RoomId: roomId, UserId: userId, Username: username}

	err := validator.New().Struct(userRoom)
	if err != nil {
		http.Error(w, "Join Room Failed at Struct Level: "+err.Error(), http.StatusBadRequest)
		return
	}

	room, exists := rr.hub.Rooms[userRoom.RoomId]
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	err = ws.ConnectToRoom(room, userRoom.UserId, userRoom.Username, rr.messageRepository, w, r)
	if err != nil {
		http.Error(w, "Failed at Connect To Room", http.StatusInternalServerError)
		return
	}
}

func (rr *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to Read Body", http.StatusInternalServerError)
		return
	}
	createRoom := &CreateRoom{}
	err = json.Unmarshal(body, createRoom)
	if err != nil {
		http.Error(w, "Failed to Unmarshal Body", http.StatusInternalServerError)
		return
	}

	err = validator.New().Struct(createRoom)
	if err != nil {
		http.Error(w, "Create Room Failed at Struct Level: "+err.Error(), http.StatusBadRequest)
		return
	}

	room := &repo.Room{Name: createRoom.RoomName}
	err = rr.roomRepository.CreateRoom(room)
	if err != nil {
		http.Error(w, "Failed to Create Room", http.StatusInternalServerError)
		return
	}

	err = rr.hub.CreateRoom(room.RoomId)
	if err != nil {
		http.Error(w, "Failed to Create Room", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

func (rr *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := rr.roomRepository.GetRooms()
	if err != nil {
		http.Error(w, "Failed to Get Rooms", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}
