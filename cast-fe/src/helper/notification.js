import React from "react";
import { authManager } from "./auth";
import { toast } from "react-toastify";
import Toast from "../components/Toast";
import api from "../apis/api";

let connection;
let ping = {};

export default function notification() {
  return {
    init: () => {
      if (!authManager.isAuthenticated()) return;
      window.WebSocket = window.WebSocket || window.MozWebSocket;
      if (!window.WebSocket) {
        console.log("WebSocket not supported!");
        return;
      }
      connection = new WebSocket(api.ws.notification());
      connection.onopen = () => {
        console.log("Notification connected!");
        ping = setInterval(
          () =>
            connection.send(
              JSON.stringify({
                type: "ping",
              })
            ),
          30000
        );
      };
      connection.onerror = () => {
        console.log("Cannot connect to notification server!");
        clearInterval(ping);
      };
      connection.onclose = () => {
        console.log("Notification server disconnected!");
        clearInterval(ping);
      };
      connection.onmessage = (message) => {
        try {
          let json = JSON.parse(message.data);
          if (json && json.type === "notification") {
            toast(<Toast data={json.data} />);
          } else if (json && json.type !== "ping") {
            console.log("Invalid JSON: ", message.data);
          }
        } catch (e) {
          console.log("Invalid JSON: ", message.data);
        }
      };
    },
    destroy: () => {
      connection.close(1000);
      clearInterval(ping);
      connection = {};
    },
  };
}
