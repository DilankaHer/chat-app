package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	ID         string
}

var hubs = map[string]*Hub{}

func GetHubs() map[string]*Hub {
	return hubs
}

func NewHub(id string) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		ID:         id,
	}
}

// func AddUserToHub(hub *Hub, userID string) {
// 	// Logic to add user to the hub can be implemented here
// 	fmt.Printf("User %s added to hub %s\n", userID, hub.ID)
// }

func (h *Hub) Run() {
	hubs[h.ID] = h
	fmt.Println("Hub started with ID:", h.ID)
	for {
		select {
		case client := <-h.register:
			fmt.Println("Client registered", client)
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				fmt.Printf("Client unregistered")
				delete(h.clients, client)
				close(client.send)
			}
			if len(h.clients) == 0 {
				delete(hubs, h.ID)
				close(h.register)
				close(h.unregister)
				close(h.broadcast)
				return
			}

		case message := <-h.broadcast:
			// Send to all clients
			fmt.Println("all client", h.clients)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	id   string
}

// read messages from browser → broadcast to hub
func (c *Client) readPump() {
	defer func() {
		fmt.Println("readPump exiting for client:", c.conn.RemoteAddr())
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Printf("Received message: %s,%s\n", c.id, message)
		var msg = Message{
			ID:      c.id,
			Content: string(message),
		}
		msgByte, err := json.Marshal(msg)
		if err != nil {
			log.Println("failed to marshal", err)
			break
		}
		c.hub.broadcast <- msgByte
	}
}

// write messages from hub → send to browser
func (c *Client) writePump() {
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

func ReadMessage(hub *Hub, clientId string, w http.ResponseWriter, r *http.Request) {
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
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		id:   clientId,
	}

	fmt.Printf("Client Created %s\n", client.id)

	hub.register <- client

	go client.writePump()
	go client.readPump()

}
