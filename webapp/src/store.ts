import { unstable_batchedUpdates } from "react-dom";
import create from "zustand";
import connect, { ConnectionApi, ConnectionState } from "./connect";
import { Building } from "./generated/building";
import { ClientMessage } from "./generated/client_message";
import { Education } from "./generated/education";
import {
  Construction,
  GameStatePatch,
  GameStatePatch_LotPatch,
  GameState_Population,
  Training,
} from "./generated/gamestate";
import { GameData } from "./generated/game_data";
import { ServerMessage, ServerMessage_User } from "./generated/server_message";

export type Lot = {
  building: Building;
};

export type GameState = {
  resources: {
    pizzas: number;
    coins: number;
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
  gameData: GameData | null;
  user: User | null;
  connection: ConnectionApi | null;
  connectionState: ConnectionState | null;
  setGameState: (gameState: GameState) => void;
  fetchGameData: () => Promise<void>;
  start: () => void;
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
      pizzas: 0,
      coins: 0,
    },
    lots: {},
    population: {
      uneducated: 0,
      chefs: 0,
      salesmice: 0,
      guards: 0,
      thieves: 0,
    },
    trainingQueue: [],
    constructionQueue: [],
  },
  user: null,
  connection: null,
  connectionState: null,
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
  start: () => {
    get().connection?.close();

    const handleStateChange = (stateChange: GameStatePatch) => {
      const resources: Partial<State["gameState"]["resources"]> = {};
      if (stateChange.resources?.coins?.value) {
        resources.coins = stateChange.resources.coins.value;
      }
      if (stateChange.resources?.pizzas?.value) {
        resources.pizzas = stateChange.resources.pizzas.value;
      }

      const population: Partial<State["gameState"]["population"]> = {};
      if (stateChange.population?.uneducated) {
        population.uneducated = stateChange.population.uneducated.value;
      }
      if (stateChange.population?.chefs) {
        population.chefs = stateChange.population.chefs.value;
      }
      if (stateChange.population?.salesmice) {
        population.salesmice = stateChange.population.salesmice.value;
      }
      if (stateChange.population?.guards) {
        population.guards = stateChange.population.guards.value;
      }
      if (stateChange.population?.thieves) {
        population.thieves = stateChange.population.thieves.value;
      }

      unstable_batchedUpdates(() => {
        set((state) => ({
          ...state,
          gameState: {
            ...state.gameState,
            resources: {
              ...state.gameState.resources,
              ...resources,
            },
            lots: mergeLots(state.gameState.lots, stateChange.lots),
            population: {
              ...state.gameState.population,
              ...population,
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
    };

    const handleUserMessage = (msg: ServerMessage_User) => {
      unstable_batchedUpdates(() => {
        set((state) => ({
          user: {
            ...state.user,
            username: msg.username,
          }
        }))
      });
    };

    const onMessage = (msg: ServerMessage) => {
      switch (msg.payload.oneofKind) {
        case "stateChange":
          handleStateChange(msg.payload.stateChange);
          break;
        case "user":
          handleUserMessage(msg.payload.user);
          break;
      }
    };
    const onStateChange = (connectionState: ConnectionState) => {
      unstable_batchedUpdates(() => {
        set((state) => ({ ...state, connectionState }));
      });
    };
    const connection = connect(onStateChange, onMessage);
    set((state) => ({ ...state, connection }));
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
