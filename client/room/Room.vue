<script setup lang="ts">
import { computed } from '@vue/reactivity';
import { watch, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const roomName = computed(() => route.query.name)

const NoRoomAction:string = "no-room"
const RoomExistsAction:string = "room-exists"
const RoomCreatedAction:string = "room-created"
const JoinedRoomAction:string = "joined-room"
const ClientConnectedAction:string = "client-connected"
const SendMessageAction:string = "send-message"
const UserJoinAction:string = "user-join"

const messages = ref([])
const messageInput = ref('')
const socket = ref()

const route = useRoute()
console.log(route.query.name);

socket.value = new WebSocket("ws://localhost:8080/ws" + "?name=" + "room")

watch(roomName, newRoomName => {
  socket.value = new WebSocket("ws://localhost:8080/ws" + "?name=" + newRoomName)


  socket.value.onopen = function() {
    socket.value.send(JSON.stringify({
      action: 'join-room',
      message: roomName.value,
      target: roomName.value,
    }));
  };

  socket.value.onmessage = function (e:MessageEvent) {
    const response = JSON.parse(e.data);
    console.log(response)
    if (response.action ==  SendMessageAction) {
      messages.value.push(response.message)
    } else if (response.action == UserJoinAction) {
      console.log("User joined")
    }
  }
})


function sendMessage() {
  socket.value.send(JSON.stringify({
    action: 'send-message',
    message: messageInput.value,
    target: roomName.value,
  }));
}

</script>

<template>
  <header>
  </header>

  <main>
    <div class="greetings">
      Room {{ roomName }}
    </div>
    <div
      v-for="message in messages"
    >
      {{ message }}
    </div>
    <input v-model="messageInput" /> 
    <button @click="sendMessage">Send</button>
  </main>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  position: relative;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h3 {
    text-align: left;
  }
}
</style>
