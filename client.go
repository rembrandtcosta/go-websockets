package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if len(origin) < 4 {
			return false
		}
		origin = origin[len(origin)-4:]
		return origin == "5173"
	},
}

type Client struct {
	ClientId uuid.UUID `json:"client-id"`
	hub *Hub
	conn *websocket.Conn
	send chan []byte
	rooms map[*Room]bool
}

func newClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ClientId: uuid.New(),
		hub: hub,
		conn: conn,
		send: make(chan []byte),
		rooms: make(map[*Room]bool),
	}
}

func (client *Client) GetName() string {
	return client.ClientId.String() 
}

func (client *Client) disconnect() {
	client.hub.unregister <- client 
	for room := range client.rooms {
		room.unregister <- client
	}
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	log.Println("ccccc")
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}
	log.Println(message)

	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		roomName := message.Target
		if room := client.hub.findRoomByName(roomName); room != nil {
			log.Println(message)
			room.broadcast <- &message 
		}

	
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case CreateRoomActon:
		client.handleCreateRoomMessage(message)
	}
}

const roomAlreadyExists = "Room already exists"
const roomDoesntExist = "Room doesn't exist"

func (client *Client) handleCreateRoomMessage(message Message) {
	roomName := message.Message

	room := client.hub.findRoomByName(roomName)
	if room != nil {
		msg := &Message{
			Action: RoomExistsAction,
			Message: roomAlreadyExists,
		}
		message.Sender.send <- msg.encode()
		return
	}

	room = client.hub.createRoom(roomName)

	client.rooms[room] = true 

	room.register <- client

	msg := &Message{
		Action: RoomCreatedAction,
		Message: roomName,
	}
	message.Sender.send <- msg.encode()
}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message
	log.Println(client.ClientId)
	log.Println(roomName)

	room := client.hub.findRoomByName(roomName)
	if room == nil {
		msg := &Message{
			Action: NoRoomAction,
			Message: roomDoesntExist,
			Sender: client,
		}
		message.Sender.send <- msg.encode()
		return
	}

	client.rooms[room] = true
	room.register <- client 
  
	msg := &Message{
		Action: JoinRoomAction,
		Message: "Joined room successfully",
		Sender: client,
	}
	message.Sender.send <- msg.encode()
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.hub.findRoomByName(message.Message)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client 
}


var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

/*
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleNewMessage(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)	
			if err != nil {
				return
			}
			log.Println(string(message))
			w.Write(message)	

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}
*/

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url param 'name' is missing")
		return 
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, hub) //&Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	go client.writePump()
	go client.readPump()

	hub.register <- client
}
