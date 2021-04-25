import { unstable_batchedUpdates } from "react-dom";
import create from "zustand";
import connect, { Connection } from "./connect";
import { ClientMessage } from "./generated/client_message";

export type Lot = {
  building: string
}

export type GameState = {
  resources: {
    pizzas: bigint;
    coins: bigint;
  };
  lots: Record<string, Lot|undefined>;
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
  constructBuilding: (lotId: string, building: string) => void;
};

const mergeLots = (a: GameState["lots"], b: Record<string, GameStatePatch_LotPatch>): GameState["lots"] => {
  const res = { ...a };

  Object.keys(b).forEach(lotId => {
    if (b[lotId] !== undefined) {
      res[lotId] = b[lotId];
    }
  });

  return res;
}

export const useStore = create<State>((set, get) => ({
  gameState: {
    resources: {
      pizzas: BigInt(0),
      coins: BigInt(0),
    },
    lots: {},
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
              lots: mergeLots(state.gameState.lots, stateChange.lots),
            },
          }));
        });
      }
    });
    set((state) => ({ ...state, user: { username }, connection }));
  },
  setGameState: (gameState) => set((state) => ({ ...state, gameState })),
  constructBuilding: (lotId, building) => {
    get().connection?.send(
      ClientMessage.create({
        id: "test-456",
        type: {
          oneofKind: "constructBuilding",
          constructBuilding: {
            building,
            lotId,
          },
        },
      })
    );
  },
}));

useStore.subscribe(console.log);
