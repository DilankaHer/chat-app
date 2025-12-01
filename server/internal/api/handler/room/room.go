package handler

import (
	"duhchat/internal/repo"
	"duhchat/internal/ws"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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
	err = ws.ConnectToRoom(room, userRoom.UserId, userRoom.Username, rr.messageRepo, w, r)
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
