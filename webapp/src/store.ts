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
  OngoingResearch,
} from "./generated/gamestate";
import { GameData } from "./generated/game_data";
import { Report } from "./generated/report";
import { ResearchDiscovery } from "./generated/research";
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
  level: number;
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
  discoveries: Array<ResearchDiscovery>;
  researchQueue: Array<OngoingResearch>;
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
  upgradeBuilding: (lotId: string) => void;
  razeBuilding: (lotId: string) => void;
  cancelRazeBuilding: (lotId: string) => void;
  train: (education: Education, amount: number) => void;
  steal: (x: number, y: number, amount: number) => void;
  readReport: (id: string) => void;
  startResearch: (discovery: ResearchDiscovery) => void;
};

const mergeLots = (
  a: GameState["lots"],
  b: Record<string, GameStatePatch_LotPatch>
): GameState["lots"] => {
  const res = { ...a };

  Object.keys(b).forEach((lotId) => {
    const lot = res[lotId];
    if (b[lotId] !== undefined) {
      const { building, tappedAt, level, razed } = b[lotId];

      if (razed) {
        delete res[lotId];
        return;
      }

      if (lot === undefined) {
        res[lotId] = { building, tappedAt, level };
      } else {
        lot.building = building;
        lot.tappedAt = tappedAt;
        lot.level = level;
      }
    }
  });

  return res;
};

const initialGameState: GameState = {
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
    publicists: 0,
  },
  trainingQueue: [],
  constructionQueue: [],
  travelQueue: [],
  townX: 0,
  townY: 0,
  discoveries: [],
  researchQueue: [],
};

const resetAuthState = (state: State) => ({
  ...state,
  connection: null,
  connectionState: null,
  user: null,
  gameState: initialGameState,
  gameStats: null,
  reports: [],
});

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
    const gameData = GameData.fromJson(data);
    set((state) => ({ ...state, gameData }));
  },
  logout: async () => {
    set(resetAuthState);
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
      if (stateChange.resources?.coins?.value !== undefined) {
        resources.coins = stateChange.resources.coins.value;
      }
      if (stateChange.resources?.pizzas?.value !== undefined) {
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
      if (stateChange.population?.publicists) {
        population.publicists = stateChange.population.publicists.value;
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
            discoveries: stateChange.discoveriesPatched
              ? stateChange.discoveries
              : state.gameState.discoveries,
            researchQueue: stateChange.researchQueuePatched
              ? stateChange.researchQueue
              : state.gameState.researchQueue,
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
        set((state) => {
          if (connectionState.error === 'unauthorized') {
            state = resetAuthState(state);
          }
          return { ...state, connectionState }
        });
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
  upgradeBuilding: (lotId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "upgradeBuilding",
          upgradeBuilding: {
            lotId,
          },
        },
      })
    );
  },
  razeBuilding: (lotId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "razeBuilding",
          razeBuilding: {
            lotId,
          },
        },
      })
    );
  },
  cancelRazeBuilding: (lotId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "cancelRazeBuilding",
          cancelRazeBuilding: {
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
            education,
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
  startResearch: (discovery) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "startResearch",
          startResearch: {
            discovery,
          }
        },
      })
    );
  }
}));

(window as any).useStore = useStore;
