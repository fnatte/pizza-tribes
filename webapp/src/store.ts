import { unstable_batchedUpdates } from "react-dom";
import create from "zustand";
import connect, { Connection } from "./connect";
import { Building } from "./generated/building";
import { ClientMessage } from "./generated/client_message";
import { Education } from "./generated/education";
import {
  Construction,
  GameStatePatch_LotPatch,
  GameState_Population,
  Training,
} from "./generated/gamestate";
import { GameData } from "./generated/game_data";

export type Lot = {
  building: Building;
};

export type GameState = {
  resources: {
    pizzas: bigint;
    coins: bigint;
  };
  lots: Record<string, Lot | undefined>;
  population: GameState_Population;
  trainingQueue: Array<Training>;
  constructionQueue: Array<Construction>;
};

type User = {
  username: string;
};

type State = {
  gameState: GameState;
  gameData: GameData|null;
  user: User | null;
  connection: Connection | null;
  setGameState: (gameState: GameState) => void;
  fetchGameData: () => Promise<void>;
  start: (username: string) => void;
  logout: () => Promise<void>;
  tap: () => void;
  constructBuilding: (lotId: string, building: Building) => void;
  train: (education: Education, amount: number) => void;
};

const mergeLots = (
  a: GameState["lots"],
  b: Record<string, GameStatePatch_LotPatch>
): GameState["lots"] => {
  const res = { ...a };

  Object.keys(b).forEach((lotId) => {
    const lot = res[lotId];
    if (b[lotId] !== undefined) {
      const { building } = b[lotId];
      if (building !== undefined) {
        if (lot === undefined) {
          res[lotId] = { building };
        } else {
          lot.building = building;
        }
      }
    }
  });

  return res;
};

export const useStore = create<State>((set, get) => ({
  gameState: {
    resources: {
      pizzas: BigInt(0),
      coins: BigInt(0),
    },
    lots: {},
    population: {
      uneducated: BigInt(0),
      chefs: BigInt(0),
      salesmice: BigInt(0),
      guards: BigInt(0),
      thieves: BigInt(0),
    },
    trainingQueue: [],
    constructionQueue: [],
  },
  user: null,
  connection: null,
  gameData: null,
  fetchGameData: async () => {
    const response = await fetch("/api/gamedata");
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      console.error("Failed to get game data");
      return;
    }
    const data = await response.json();
    const gameData = data as GameData;
    set((state) => ({ ...state, gameData }));
  },
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
              population: {
                ...state.gameState.population,
                ...stateChange.population,
              },
              trainingQueue: stateChange.trainingQueuePatched
                ? stateChange.trainingQueue
                : state.gameState.trainingQueue,
              constructionQueue: stateChange.constructionQueuePatched
                ? stateChange.constructionQueue
                : state.gameState.constructionQueue,
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
  train: (education, amount) => {
    get().connection?.send(
      ClientMessage.create({
        id: "test-456",
        type: {
          oneofKind: "train",
          train: {
            education: education,
            amount,
          },
        },
      })
    );
  },
}));
