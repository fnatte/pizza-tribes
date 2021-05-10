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
  Travel,
} from "./generated/gamestate";
import { GameData } from "./generated/game_data";
import { Report } from "./generated/report";
import {
  ServerMessage,
  ServerMessage_Reports,
  ServerMessage_User,
} from "./generated/server_message";
import { Stats } from "./generated/stats";
import { generateId } from "./utils";

export type Lot = {
  building: Building;
  tappedAt: string;
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
  travelQueue: Array<Travel>;
  townX: number;
  townY: number;
};

type User = {
  username: string;
};

type State = {
  gameState: GameState;
  gameStats: Stats | null;
  gameData: GameData | null;
  gameDataLoading: boolean;
  reports: Report[];
  user: User | null;
  connection: ConnectionApi | null;
  connectionState: ConnectionState | null;
  setGameState: (gameState: GameState) => void;
  fetchGameData: () => Promise<void>;
  start: () => void;
  logout: () => Promise<void>;
  tap: (lotId: string) => void;
  constructBuilding: (lotId: string, building: Building) => void;
  train: (education: Education, amount: number) => void;
  steal: (x: number, y: number, amount: number) => void;
  readReport: (id: string) => void;
};

const mergeLots = (
  a: GameState["lots"],
  b: Record<string, GameStatePatch_LotPatch>
): GameState["lots"] => {
  const res = { ...a };

  Object.keys(b).forEach((lotId) => {
    const lot = res[lotId];
    if (b[lotId] !== undefined) {
      const { building, tappedAt } = b[lotId];
      if (lot === undefined) {
        res[lotId] = { building, tappedAt };
      } else {
        lot.building = building;
        lot.tappedAt = tappedAt;
      }
    }
  });

  return res;
};

const initialGameState = {
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
  travelQueue: [],
  townX: 0,
  townY: 0,
};

export const useStore = create<State>((set, get) => ({
  gameState: initialGameState,
  gameStats: null,
  user: null,
  connection: null,
  connectionState: null,
  gameData: null,
  gameDataLoading: false,
  reports: [],
  fetchGameData: async () => {
    set((state) => ({ ...state, gameDataLoading: true }));
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
    set((state) => ({
      ...state,
      connection: null,
      connectionState: null,
      user: null,
      gameState: initialGameState,
      gameStats: null,
      reports: [],
    }));
    const res = await fetch("/api/auth/logout");
    if (res.ok) {
      get().connection?.close();
    }
  },
  tap: (lotId: string) => {
    const msg = ClientMessage.create({
      id: generateId(),
      type: {
        oneofKind: "tap",
        tap: { lotId },
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
            travelQueue: stateChange.travelQueuePatched
              ? stateChange.travelQueue
              : state.gameState.travelQueue,
            townX: stateChange.townX?.value ?? state.gameState.townX,
            townY: stateChange.townY?.value ?? state.gameState.townY,
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
          },
        }));
      });
    };

    const handleStats = (stats: Stats) => {
      unstable_batchedUpdates(() => {
        set((state) => ({
          ...state,
          gameStats: stats,
        }));
      });
    };

    const handleReports = (msg: ServerMessage_Reports) => {
      unstable_batchedUpdates(() => {
        set((state) => ({
          ...state,
          reports: msg.reports,
        }));
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
        case "stats":
          handleStats(msg.payload.stats);
          break;
        case "reports":
          handleReports(msg.payload.reports);
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
        id: generateId(),
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
        id: generateId(),
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
  steal: (x: number, y: number, amount: number) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "steal",
          steal: {
            x,
            y,
            amount,
          },
        },
      })
    );
  },
  readReport: (id: string) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "readReport",
          readReport: { id },
        },
      })
    );

    set((state) => {
      const reports = [...state.reports];
      const i = reports.findIndex((report) => report.id === id);
      if (i !== -1) {
        reports[i] = { ...reports[i], unread: false };
      }
      return { ...state, reports };
    });
  },
}));

(window as any).useStore = useStore;
