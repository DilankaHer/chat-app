package ws

import (
	"duhchat/internal/api/model"
	"duhchat/internal/repo"
	"fmt"
	"log"
)

type Hub struct {
	Rooms          map[string]Room
	roomRepository repo.RoomRepository
}

type Room struct {
	users      map[*User]bool
	broadcast  chan []byte
	register   chan *User
	unregister chan *User
	RoomId     string
}

func NewHub(roomRepository repo.RoomRepository) *Hub {
	return &Hub{
		Rooms:          make(map[string]Room),
		roomRepository: roomRepository,
	}
}

func (h *Hub) CreateRoom(roomId string) error {
	h.Rooms[roomId] = Room{
		users:      make(map[*User]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *User),
		unregister: make(chan *User),
		RoomId:     roomId,
	}

	fmt.Println("New Room ID", roomId)
	err := h.AddNewRoom(roomId)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hub) CreateDefaultRooms() error {
	roomIds, err := h.roomRepository.GetRoomIds()
	if err != nil {
		log.Fatal("Failed to get room IDs:", err)
		return err
	}
	for _, roomId := range roomIds {
		h.CreateRoom(roomId)
	}
	return nil
}

func (h *Hub) AddNewRoom(roomId string) error {
	room := h.Rooms[roomId]
	go func() error {
		for {
			select {
			case user := <-room.register:
				fmt.Println("User registered", user)
				room.users[user] = true
				err := h.roomRepository.JoinRoom(&model.UserRoom{
					UserId: user.userId,
					RoomId: room.RoomId,
				})
				if err != nil {
					delete(room.users, user)
					close(user.send)
					log.Println("Failed to join room:", err)
					return err
				}
			case user := <-room.unregister:
				if _, ok := room.users[user]; ok {
					fmt.Printf("User unregistered")
					delete(room.users, user)
					close(user.send)
					err := h.roomRepository.DeleteRoomUsersByUserId(user.userId)
					if err != nil {
						log.Println("Failed to delete room users:", err)
						return err
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

	return nil
}
