package message

import (
	"duhchat/internal/repo"
	"encoding/json"
	"net/http"
)

type MessageHandler struct {
	messageRepo repo.MessageRepository
}

func NewMessageHandler(messageRepo repo.MessageRepository) *MessageHandler {
	return &MessageHandler{messageRepo: messageRepo}
}

func (mh *MessageHandler) GetMessagesByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	messages, err := mh.messageRepo.GetMessagesByRoomId(roomId)
	if err != nil {
		http.Error(w, "Failed to Get Messages By Room ID", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
