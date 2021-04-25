import {ClientMessage} from "./generated/client_message";
import {ServerMessage} from "./generated/server_message";

export type Connection = {
  send: (msg: ClientMessage) => void;
  close: () => void;
};

type MessageHandler = (msg: ServerMessage) => void;

const connect = (onMessage: MessageHandler): Connection => {
  const conn = new WebSocket("ws://localhost:3000/api/ws");

  const send = (msg: ClientMessage) => {
    conn.send(JSON.stringify(ClientMessage.toJson(msg)));
  };

  conn.onclose = (e) => {
    const unauthorized = e.code === 4010;
    console.log("unauthorized?", unauthorized);
  };

  conn.onmessage = (e) => {
    const message = ServerMessage.fromJson(JSON.parse(e.data))
    onMessage(message);
  };

  conn.onopen = () => {
    console.log("connected");
  }

  conn.onerror = (e) => {
    console.log(e);
  };

  return {
    send,
    close: () => {
      conn.close();
      console.log('closing', conn);
    },
  };
};

export default connect;

