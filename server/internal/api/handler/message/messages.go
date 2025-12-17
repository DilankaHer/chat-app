package message

import (
	"duhchat/internal/repo"
	"duhchat/util"
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
		util.JSONMarshaller(w, http.StatusInternalServerError, "Failed to Get Messages", http.StatusText(http.StatusInternalServerError))
		return
	}
	util.JSONMarshaller(w, http.StatusOK, messages, http.StatusText(http.StatusOK))
}
