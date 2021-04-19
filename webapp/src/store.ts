import { unstable_batchedUpdates } from "react-dom";
import create, {
  GetState,
  SetState,
  StateCreator,
  StoreApi,
  UseStore,
} from "zustand";
import connect, { Connection } from "./connect";

type GameState = {
  resources: {
    pizzas: number;
    coins: number;
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
  start: () => void;
  logout: () => Promise<void>;
};

export const useStore = create<State>((set, get) => ({
  gameState: {
    resources: {
      pizzas: 0,
      coins: 0,
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
  start: () => {
    const connection = connect((msg) => {
      unstable_batchedUpdates(() => {
        set((state) => ({
          ...state,
          gameState: { ...state.gameState, ...msg },
        }));
      });
    });
    set((state) => ({ ...state, user: { username: "" }, connection }));
  },
  setGameState: (gameState) => set((state) => ({ ...state, gameState })),
}));

useStore.subscribe(console.log);

