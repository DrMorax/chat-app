export const socketHandler = {
  socketEvents: {
    handleOpen: function () {
      console.log("Connected to the server");
    },

    handleError: function (error) {
      console.error("WebSocket error:", error);
    },

    handleMessage: function (event) {
      const message = JSON.parse(event.data);
      console.log(`Received message: ${message.text}`);
      document.getElementById(
        "messages"
      ).innerHTML += `<p>From server: ${message.text}</p>`;
    },

    handleClose: function () {
      console.log("Disconnected from the server");
    },
  },

  eventListeners: {
    clickSend: function (socket) {
      document
        .getElementById("send-message-btn")
        .addEventListener("click", () => {
          this.actions.sendMessage(socket);
        });
    },
    enterSend: function (socket) {
      document
        .getElementById("message-input")
        .addEventListener("keyup", (event) => {
          if (event.key === "Enter") {
            this.actions.sendMessage(socket);
          }
        });
    },
    clickDisconnect: function (socket) {
      document
        .getElementById("disconnect-btn")
        .addEventListener("click", () => {
          this.actions.closeConnection(socket);
        });
    },
    actions: {
      sendMessage: function (socket) {
        let input = document.getElementById("message-input");
        if (input.value === "") {
          console.error("Message cannot be empty");
          return;
        }
        socket.send(input.value);
        document.getElementById(
          "messages"
        ).innerHTML += `<p>From client: ${input.value}</p>`;
        input.value = "";
      },
      closeConnection: function (socket) {
        socket.close();
      },
    },
  },
};
