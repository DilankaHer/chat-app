package handler

import (
	"duhchat/internal/repo"
	"duhchat/internal/ws"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type JoinRoomHandler struct {
	joinRoomRepo repo.JoinRoomRepository
	messageRepo  repo.MessageRepository
	hub          *ws.Hub
}

func NewJoinRoomHandler(joinRoomRepo repo.JoinRoomRepository, messageRepo repo.MessageRepository, hub *ws.Hub) *JoinRoomHandler {
	return &JoinRoomHandler{joinRoomRepo: joinRoomRepo, messageRepo: messageRepo, hub: hub}
}

func (rr *JoinRoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	userId := r.URL.Query().Get("userId")

	userRoom := &repo.UserRoom{RoomId: roomId, UserId: userId}

	err := validator.New().Struct(userRoom)
	if err != nil {
		http.Error(w, "Join Room Failed at Struct Level: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Available Hub IDs:", rr.hub.Rooms)

	room, exists := rr.hub.Rooms[userRoom.RoomId]
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	err = ws.ConnectToRoom(room, userRoom.UserId, rr.messageRepo, w, r)
	if err != nil {
		http.Error(w, "Failed at Connect To Room", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Join Room Successful"))
}

func (rr *JoinRoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := rr.joinRoomRepo.GetRooms()
	if err != nil {
		http.Error(w, "Failed to Get Rooms", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}
