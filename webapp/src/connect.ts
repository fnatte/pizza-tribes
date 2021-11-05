import {WS_URL} from "./config";
import { ClientMessage } from "./generated/client_message";
import { ServerMessage } from "./generated/server_message";

type MessageHandler = (msg: ServerMessage) => void;

export type ConnectionState = {
  connected: boolean;
  connecting: boolean;
  reconnectAttempts: number;
  error: "unknown" | "unauthorized" | false;
};

export type ConnectionApi = {
  reconnect: () => void;
  send: (msg: ClientMessage) => void;
  close: () => void;
};

const connect = (
  onStateChange: (state: ConnectionState) => void,
  onMessage: MessageHandler
): ConnectionApi => {
  let conn: WebSocket | null = null;
  let state: ConnectionState = {
    connected: false,
    connecting: false,
    reconnectAttempts: 0,
    error: false,
  };
  let targetState: "connected" | "disconnected" = "connected";
  let pendingReconnectAttempt: number | null = null;

  const setState = (p: Partial<ConnectionState>): void => {
    state = { ...state, ...p };
    onStateChange(state);
  };

  const makeReconnectAttempt = () => {
    if (pendingReconnectAttempt !== null) {
      return;
    }

    if (state.reconnectAttempts >= 100) {
      targetState = "disconnected";
      setState({ connecting: false });
      return;
    }

    if (targetState === "disconnected") {
      return;
    }

    setState({ reconnectAttempts: state.reconnectAttempts + 1 });

    const delay = Math.min(
      state.reconnectAttempts * state.reconnectAttempts * 1000,
      5_000
    );

    pendingReconnectAttempt = window.setTimeout(() => {
      pendingReconnectAttempt = null;
      reconnect();
    }, delay);
  };

  const reconnect = () => {
    if (targetState === "disconnected") {
      return;
    }

    if (pendingReconnectAttempt !== null) {
      window.clearTimeout(pendingReconnectAttempt);
      pendingReconnectAttempt = null;
    }

    setState({ connecting: true });

    conn?.close();
    conn = new WebSocket(WS_URL);
    conn.onclose = (e) => {
      const initializationError = e.code === 5001;
      const unauthorized = e.code === 4010;
      if (unauthorized || initializationError) {
        targetState = "disconnected";
        setState({
          error: "unauthorized",
          connecting: false,
          connected: false,
          reconnectAttempts: 0,
        });
        return;
      }

      if (targetState === "connected") {
        makeReconnectAttempt();
      }
    };

    conn.onmessage = (e) => {
      const message = ServerMessage.fromJson(JSON.parse(e.data));
      onMessage(message);
    };

    conn.onopen = () => {
      setState({ connected: true, connecting: false, reconnectAttempts: 0 });
    };

    conn.onerror = (e) => {
      console.error("WebSocket error", e);
    };
  };

  const send = (msg: ClientMessage) => {
    conn?.send(JSON.stringify(ClientMessage.toJson(msg)));
  };

  reconnect();

  return {
    send,
    reconnect: () => {
      targetState = "connected";
      if (
        conn === null ||
        conn.readyState === WebSocket.CLOSED ||
        conn.readyState === WebSocket.CLOSING
      ) {
        reconnect();
      }
    },
    close: () => {
      targetState = "disconnected";
      if (pendingReconnectAttempt !== null) {
        window.clearTimeout(pendingReconnectAttempt);
        pendingReconnectAttempt = null;
      }
      if (state.connecting) {
        setState({ connecting: false });
      }
      conn?.close();
      conn = null;
    },
  };
};

export default connect;
