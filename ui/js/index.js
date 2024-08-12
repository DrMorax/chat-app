import { socketHandler } from "./socket.js";

const socket = new WebSocket("ws://localhost:3001/ws/123?v=1.0");

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
