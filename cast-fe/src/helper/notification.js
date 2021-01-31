import React from "react";
import urls from "./url";
import auth from "./auth";
import { toast } from "react-toastify";
import Toast from "../components/Toast";

let connection;
let ping = {};

export default function notification() {
  return {
    init: () => {
      if (!auth().is_authenticated()) return;
      window.WebSocket = window.WebSocket || window.MozWebSocket;
      if (!window.WebSocket) {
        console.log("WebSocket not supported!");
        return;
      }
      connection = new WebSocket(urls().notification(auth().token()));
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
