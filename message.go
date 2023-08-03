package main

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const CreateRoomActon = "create-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"

const NoRoomAction = "no-room"
const RoomExistsAction = "room-exists"
const RoomCreatedAction = "room-created"
const JoinedRoomAction = "joined-room"
const ClientConnectedAction = "client-connected"

type Message struct {
	Action string `json:"action"`
	Message string `json:"message"`
	Target string `json:"target"`
	Sender *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json 
}

