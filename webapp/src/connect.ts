import { getAccessToken } from "./api";
import { WS_URL } from "./config";
import { ClientMessage } from "./generated/client_message";
import { ServerMessage } from "./generated/server_message";

type MessageHandler = (msg: ServerMessage) => void;

export type ConnectionState = {
  connected: boolean;
  connecting: boolean;
  reconnectAttempts: number;
  error: "unknown" | "unauthorized" | "reconnect-failed" | false;
};

export type ConnectionApi = {
  reconnect: () => void;
  send: (msg: ClientMessage) => void;
  close: () => void;
  reset: () => void;
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
      setState({ connecting: false, error: "reconnect-failed" });
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

    if (state.connecting) {
      return
    }

    if (pendingReconnectAttempt !== null) {
      window.clearTimeout(pendingReconnectAttempt);
      pendingReconnectAttempt = null;
    }

    setState({ connecting: true });

    // Wait what? Why are we setting the websocket protocol to "accessToken.${accessToken}"?
    // Yes, it's a dirty hack, but we live with it to avoid a complex authorization flow for
    // web sockets. In short, "Sec-WebSocket-Protocol" is the only header we can send stuff via
    // because web sockets do not allow custom headers such as "Authorization".
    //
    // Also, note that this is only used when we do not rely on cookie authorization, that is,
    // we only send the access token over Sec-WebSocket-Protocol for cross-origin usages (as of now,
    // mobile apps).
    const accessToken = getAccessToken();
    const protocols = accessToken
      ? ["pizzatribes", `accessToken.${accessToken}`]
      : "pizzatribes";

    const handleClose = (code: number) => {
      const initializationError = code === 5001;
      const unauthorized = code === 4010;
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

    conn?.close();
    conn = new WebSocket(WS_URL, protocols);
    conn.onclose = (e) => {
      handleClose(e.code);
    };

    conn.onmessage = (e) => {
      const json = JSON.parse(e.data);
      if (json.type === "control" && typeof json.control === "object") {
        if (json.control.type === "close") {
          const code =
            typeof json.control.code === "number" ? json.control.code : -1;
          conn?.close();
          handleClose(code);
        }
        return;
      }
      const message = ServerMessage.fromJson(json);
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
      setState({ connected: false });
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
    reset: () => {
      if (state.connected || state.connecting) {
        throw new Error(
          "Cannot reset connect state while connected or connecting"
        );
      }

      setState({
        connected: false,
        connecting: false,
        reconnectAttempts: 0,
        error: false,
      });
    },
  };
};

export default connect;
