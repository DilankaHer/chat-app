package handler

import (
	"duhchat/internal/api/model"
	"duhchat/internal/repo"
	"duhchat/internal/ws"
	"duhchat/util"
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
	if roomId == "" {
		util.JSONMarshaller(w, http.StatusBadRequest, "Room ID is required", http.StatusText(http.StatusBadRequest))
		return
	}

	claims := r.Context().Value("user").(jwt.MapClaims)
	if claims == nil {
		util.JSONMarshaller(w, http.StatusInternalServerError, "Can't find userId ", http.StatusText(http.StatusInternalServerError))
		return
	}
	userId := claims["userId"].(string)
	username := claims["username"].(string)

	userRoom := &model.UserRoom{RoomId: roomId, UserId: userId, Username: username}

	err := validator.New().Struct(userRoom)
	if err != nil {
		util.JSONMarshaller(w, http.StatusBadRequest, "Join Room Failed at Struct Level: "+err.Error(), http.StatusText(http.StatusBadRequest))
		return
	}

	room, exists := rr.hub.Rooms[userRoom.RoomId]
	if !exists {
		util.JSONMarshaller(w, http.StatusNotFound, "Room not found", http.StatusText(http.StatusNotFound))
		return
	}
	err = ws.ConnectToRoom(room, userRoom.UserId, userRoom.Username, rr.messageRepository, w, r)
	if err != nil {
		util.JSONMarshaller(w, http.StatusInternalServerError, "Failed at Connect To Room", http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (rr *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(jwt.MapClaims)["userId"].(string)
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

	room := &repo.Room{Name: createRoom.RoomName, CreatedBy: userId}
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
	util.JSONMarshaller(w, http.StatusOK, room, "Room created successfully")
}

func (rr *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := rr.roomRepository.GetRooms()
	if err != nil {
		http.Error(w, "Failed to Get Rooms", http.StatusInternalServerError)
		return
	}
	util.JSONMarshaller(w, http.StatusOK, rooms, http.StatusText(http.StatusOK))
}

func (rr *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	if roomId == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	err := rr.roomRepository.DeleteRoom(roomId)
	if err != nil {
		http.Error(w, "Failed to Delete Room", http.StatusInternalServerError)
		return
	}
	util.JSONMarshaller(w, http.StatusOK, "Room deleted successfully", "Room deleted successfully")
}
