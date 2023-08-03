package main

import (
	"log"

)

type Hub struct {
	rooms map[*Room]bool
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client

}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms: make(map[*Room]bool),
	}
}

func (h *Hub) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.RunRoom()
	h.rooms[room] = true

	return room 
}

func (h *Hub) findRoomByName(name string) *Room {
	var foundRoom *Room 
	for room := range h.rooms {
		if room.GetName() == name {
			foundRoom = room 
			break
		}
	}

	return foundRoom
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		
		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.notifyClientJoined(client)
	h.listOnlineClients(client)
	log.Println(client.GetName())
	h.clients[client] = true
	message := &Message{
		Action: ClientConnectedAction,
		Message: client.GetName(),
	}
	client.send <- message.encode()
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		h.notifyClientLeft(client)
	}
}

func (h *Hub) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client, 
	}

	h.broadcastToClients(message.encode())
}

func (h *Hub) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client, 
	}
	
	h.broadcastToClients(message.encode())
}

func (h *Hub) listOnlineClients(client *Client) {
	for existingClient := range h.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: existingClient,
		}
		client.send <- message.encode()
	}
}

func (h *Hub) broadcastToClients(message []byte) {
	log.Println("aa");
	for client := range h.clients {
		log.Println(client.GetName())
		client.send <- message 
	}
}


