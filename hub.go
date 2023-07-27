package main

import (
	"log"
	"encoding/json"

	"github.com/google/uuid"
)

type Hub struct {
	rooms map[*Room]bool

	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client

}

type Payload struct {
	Method string `json:"method"`
	ClientId uuid.UUID `json:"clientId"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) createRun(name string) *Room {
	room := newRoom()
	go room.run()
	h.rooms[room] = true

	return room 
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			payload := Payload{Method: "connect", ClientId: uuid.New()}
			ret, err := json.Marshal(payload)
			if err != nil {
				log.Println(err)
			}
			client.send <- ret
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
