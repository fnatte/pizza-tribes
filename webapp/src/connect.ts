export type Connection = {
  send: (msg: any) => void;
  close: () => void;
};

type MessageHandler = (msg: any) => void;

const connect = (onMessage: MessageHandler): Connection => {
  const conn = new WebSocket("ws://localhost:3000/api/ws");

  const send = (msg: any) => conn.send(JSON.stringify(msg));

  conn.onclose = (e) => {
    const unauthorized = e.code === 4010;
    console.log("unauthorized?", unauthorized);
  };

  conn.onmessage = (e) => {
    const message = JSON.parse(e.data);
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

