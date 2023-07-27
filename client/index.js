console.log("aaa")

var socket = new WebSocket("ws://localhost:8080/ws")
var clientId = null
var joinBtn = document.getElementById("join");
var createBtn = document.getElementById("create");

socket.onopen = function() {
  socket.send(JSON.stringify({"method": "connect"}));
};

createBtn.onclick = function(_) {
  socket.send(JSON.stringify({"method": "create"}));
}

joinBtn.onclick = function(_) {
  socket.send(JSON.stringify({"test": "test"}));
};

socket.onclose = function (_) {
  console.log("Connection closed")
}

socket.onmessage = function (e) {
  console.log(e.data);
  const response = e.data;
  if (response.method == "connect") {
    clientId = response.clientId;
    console.log(clientId);
  }
}
