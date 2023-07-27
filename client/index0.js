var socket = new WebSocket("ws://localhost:8080/ws");
var output = document.getElementById("output");
var joinBtn = document.getElementById("join");
let clientId = null;
let roomId = null;
let players = []

socket.onopen = function () {
  socket.send(JSON.stringify({"method": "connect"})) 
  output.innerHTML += "Status: Connected\n";
};

socket.onmessage = function (e) {
    const response = JSON.parse(e.data);
    if (response.method === "connect"){
      clientId = response["client-id"]
      console.log("Client id set sucessfully " + clientId)
      roomId = response["room-id"]
      console.log("Room id set sucessfully " + roomId)
    }
    if (response.method === "join"){
      players.push(response["client-id"])
      console.log("Now playing: ", players)
    }
    output.innerHTML += "\nServer: " + e.data + "\n";
};

joinBtn.onclick = function join() {
  const payload = {
    "method": "join",
    "client-id" : clientId,
    "room-id": roomId,
    "name": "Lewis",
  }
  socket.send(JSON.stringify(payload));
}
