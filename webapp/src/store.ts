import { unstable_batchedUpdates } from "react-dom";
import create from "zustand";
import connect, { Connection } from "./connect";
import { ClientMessage } from "./generated/client_message";

type GameState = {
  resources: {
    pizzas: bigint;
    coins: bigint;
  };
};

type User = {
  username: string;
};

type State = {
  gameState: GameState;
  user: User | null;
  connection: Connection | null;
  setGameState: (gameState: GameState) => void;
  start: (username: string) => void;
  logout: () => Promise<void>;
  tap: () => void;
};

export const useStore = create<State>((set, get) => ({
  gameState: {
    resources: {
      pizzas: BigInt(0),
      coins: BigInt(0),
    },
  },
  user: null,
  connection: null,
  logout: async () => {
    const res = await fetch("/api/auth/logout");
    if (res.ok) {
      get().connection?.close();
      set((state) => ({ ...state, connection: null, user: null }));
    }
  },
  tap: () => {
    const msg = ClientMessage.create({
      id: "test-123",
      type: {
        oneofKind: "tap",
        tap: {
          amount: 10,
        },
      },
    });
    get().connection?.send(msg);
  },
  start: (username: string) => {
    const connection = connect((msg) => {
      if (msg.payload.oneofKind === "stateChange") {
        const stateChange = msg.payload.stateChange;
        unstable_batchedUpdates(() => {
          set((state) => ({
            ...state,
            gameState: {
              ...state.gameState,
              resources: {
                ...state.gameState.resources,
                ...stateChange.resources,
              },
            },
          }));
        });
      }
    });
    set((state) => ({ ...state, user: { username }, connection }));
  },
  setGameState: (gameState) => set((state) => ({ ...state, gameState })),
}));

useStore.subscribe(console.log);
