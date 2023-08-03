<script setup lang="ts">

import { ref } from 'vue'

const roomLink = ref('')
const roomName = ref('room')
const currentUrl =  "ws://localhost:8080/ws"
const roomAlreadyExistsError = ref(false);
const socket = new WebSocket(currentUrl + "?name=" + roomName.value)

const NoRoomAction:string = "no-room"
const RoomExistsAction:string = "room-exists"
const RoomCreatedAction:string = "room-created"
const JoinedRoomAction:string = "joined-room"
const ClientConnectedAction:string = "client-connected"

socket.onopen = function() {
  socket.send(JSON.stringify({
    action: 'send-message',
    message: 'default',
    target: roomName.value,
  }));
};

function createRoom() {
  roomAlreadyExistsError.value = false;
  roomLink.value = roomName.value;
  socket.send(JSON.stringify({
    action: 'create-room',
    message: roomLink.value,
    target: roomLink.value,
  }))
}

socket.onmessage = function(e) {
  const response = JSON.parse(e.data);
  console.log(response);
  if (response.action == RoomExistsAction) {
    console.log("Room with name " + roomName.value + " already exists");
    roomAlreadyExistsError.value = true;
    roomLink.value = '';
  }
}

function showRoomLink() {
  return roomLink.value != ''
}

</script>

<template>
  <header>
    <h1>Create Room</h1>
  </header>

  <main>
    <input v-model="roomName" /> 
    <button @click="createRoom">Create</button>
    <div v-if="showRoomLink()">
      <h2>Room link:</h2>
      <a v-bind:href="'/room/?name='+roomLink">link</a>
    </div>
    <div v-if="roomAlreadyExistsError">
      <h2> There already exists a room with that name </h2>
    </div>
  </main>
</template>

<style scoped>
  
</style>
