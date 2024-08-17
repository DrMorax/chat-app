import { socketHandler } from "./socket.js";
import { authenticator } from "./auth.js";

document.getElementById("login-btn").addEventListener("click", () => {
  authenticator.authenticate();
});
document.getElementById("register-btn").addEventListener("click", () => {
  authenticator.register();
});

const socket = new WebSocket("ws://localhost:4000/feed");

socket.onopen = () => {
  socketHandler.socketEvents.handleOpen();
};

socket.onerror = (error) => {
  socketHandler.socketEvents.handleError(error);
};

socket.onmessage = (message) => {
  socketHandler.socketEvents.handleMessage(message);
};

socket.onclose = () => {
  socketHandler.socketEvents.handleClose();
};

socketHandler.eventListeners.clickSend(socket);
socketHandler.eventListeners.enterSend(socket);
socketHandler.eventListeners.clickDisconnect(socket);
