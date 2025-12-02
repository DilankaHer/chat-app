package message

import (
	"duhchat/internal/repo"
	"encoding/json"
	"net/http"
)

type MessageHandler struct {
	messageRepository repo.MessageRepository
}

func NewMessageHandler(messageRepository repo.MessageRepository) *MessageHandler {
	return &MessageHandler{messageRepository: messageRepository}
}

func (mh *MessageHandler) GetMessagesByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("roomId")
	messages, err := mh.messageRepository.GetMessagesByRoomId(roomId)
	if err != nil {
		http.Error(w, "Failed to Get Messages By Room ID", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
