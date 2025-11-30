package ws

import (
	"duhchat/internal/repo"
	"fmt"
	"log"
)

type Hub struct {
	Rooms        map[string]Room
	joinRoomRepo repo.JoinRoomRepository
}

type Room struct {
	users      map[*User]bool
	broadcast  chan []byte
	register   chan *User
	unregister chan *User
	ID         string
}

func NewHub(joinRoomRepo repo.JoinRoomRepository) *Hub {
	return &Hub{
		Rooms:        make(map[string]Room),
		joinRoomRepo: joinRoomRepo,
	}
}

func (h *Hub) CreateRoom(id string) {
	h.Rooms[id] = Room{
		users:      make(map[*User]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *User),
		unregister: make(chan *User),
		ID:         id,
	}
}

func (h *Hub) Run() error {
	roomIds, err := h.joinRoomRepo.GetRoomIds()
	if err != nil {
		log.Fatal("Failed to get room IDs:", err)
		return err
	}
	for _, roomId := range roomIds {
		h.CreateRoom(roomId)
		room := h.Rooms[roomId]
		go func() {
			for {
				select {
				case user := <-room.register:
					fmt.Println("User registered", user)
					room.users[user] = true
					err := h.joinRoomRepo.JoinRoom(&repo.UserRoom{
						UserId: user.id,
						RoomId: room.ID,
					})
					if err != nil {
						delete(room.users, user)
						close(user.send)
						log.Println("Failed to join room:", err)
						return
					}
				case user := <-room.unregister:
					if _, ok := room.users[user]; ok {
						fmt.Printf("User unregistered")
						delete(room.users, user)
						close(user.send)
						err := h.joinRoomRepo.DeleteRoomUsersByUserId(user.id)
						if err != nil {
							log.Println("Failed to delete room users:", err)
							return
						}
					}
				case message := <-room.broadcast:
					fmt.Println("all user", room.users)
					for user := range room.users {
						select {
						case user.send <- message:
						default:
							delete(room.users, user)
							close(user.send)
						}
					}
				}
			}
		}()
	}

	return nil
}
