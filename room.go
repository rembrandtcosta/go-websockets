package main

import (
	"fmt"
	"log"
)

type Room struct {
	name string
	clients map[*Client]bool
	broadcast chan *Message
	register chan *Client
	unregister chan *Client
}

func NewRoom(name string) *Room {
	return &Room{
		name: name,
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Room) RunRoom() {
	for {
		select {

		case client := <-h.register:
			h.registerClientInRoom(client)

		case client := <-h.unregister:
			h.unregisterClientInRoom(client)

		case message := <-h.broadcast:
			h.broadcastToClientsInRoom(message.encode())
		}
	}
}

const welcomeMessage = "%s joined the room"

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action: SendMessageAction,
		Target: room.name, 
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message.encode())
}


func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	log.Println("bbb")
	for client := range room.clients {
		log.Println(client.GetName())
		client.send <- message
	}
}

func (h *Room) GetName() string {
	return h.name
}
