package ws

import (
	"duhchat/internal/repo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	room     *Room
	conn     *websocket.Conn
	send     chan []byte
	userId   string
	username string
}

func ConnectToRoom(room Room, userId string, username string, messageRepository repo.MessageRepository, w http.ResponseWriter, r *http.Request) error {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // allow all origins
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return err
	}

	user := &User{
		room:     &room,
		conn:     conn,
		send:     make(chan []byte, 256),
		userId:   userId,
		username: username,
	}

	fmt.Printf("User Created %s\n", user.userId)

	room.register <- user

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetPongHandler(func(string) error {
		fmt.Println("pong received")
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	go user.writePump()
	go user.readPump(messageRepository)
	go user.ping()

	return nil
}

func (c *User) readPump(messageRepo repo.MessageRepository) {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var msg = repo.Message{
			RoomId:   c.room.RoomId,
			UserId:   c.userId,
			Username: c.username,
			Content:  string(message),
			IsError:  false,
		}
		msgByte, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to marshal", err)
			msg.IsError = true
			c.conn.WriteJSON(msg)
			continue
		}
		err = messageRepo.AddMessasge(&msg)
		if err != nil {
			log.Println("failed to send message", err)
			msg.IsError = true
			c.conn.WriteJSON(msg)
			continue
		}
		c.room.broadcast <- msgByte
	}
}

func (c *User) writePump() {
	defer c.conn.Close()
	for {
		msg, ok := <-c.send
		if !ok {
			fmt.Println("Channel closed, closing connection")
			return
		}
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (c *User) ping() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("ping executed")
		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			fmt.Println("ping error", err)
			return
		}
	}
}
