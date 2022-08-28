import { unstable_batchedUpdates } from "react-dom";
import create from "zustand";
import connect, { ConnectionApi, ConnectionState } from "./connect";
import { Building } from "./generated/building";
import { ClientMessage } from "./generated/client_message";
import { Education } from "./generated/education";
import { GameState, GameStatePatch } from "./generated/gamestate";
import { GameData } from "./generated/game_data";
import { Report } from "./generated/report";
import { ResearchDiscovery } from "./generated/research";
import {
  ServerMessage,
  ServerMessage_Reports,
  ServerMessage_User,
} from "./generated/server_message";
import { Stats } from "./generated/stats";
import { MouseAppearance } from "./generated/appearance";
import { generateId } from "./utils";
import { enablePatches, produce } from "immer";
import { queryClient } from "./queryClient";
import { apiFetch, centralApiFetch } from "./api";
import { extractMessage } from "./protobuf/extractMessage";
import { reflectionMergePartial } from "./protobuf/protobuf-ts/reflectionMergePartial";

enablePatches();

export type Lot = {
  building: Building;
  tappedAt: string;
  level: number;
  taps: number;
  streak: number;
};

export { GameState };

type User = {
  username: string;
};

export type State = {
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
  reschool: (mouseId: string) => void;
  renameMouse: (mouseId: string, name: string) => void;
  openQuest: (questId: string) => void;
  claimQuestReward: (questId: string) => void;
  completeQuest: (questId: string) => void;
  reportActivity: () => void;
  saveMouseAppearance: (mouseId: string, appearance: MouseAppearance) => void;
  setAmbassadorMouse: (mouseId: string) => void;
  setPizzaPrice: (pizzaPrice: number) => void;
  buyGeniusFlash: (id: number) => void;
};

const resetQueryDataState = () => {
  queryClient.setQueriesData({}, () => undefined);
};

const initialGameState: GameState = {
  resources: {
    pizzas: 0,
    coins: 0,
  },
  lots: {},
  trainingQueue: [],
  constructionQueue: [],
  travelQueue: [],
  townX: 0,
  townY: 0,
  discoveries: [],
  researchQueue: [],
  mice: {},
  quests: {},
  timestamp: "",
  ambassadorMouseId: "",
  pizzaPrice: 0,
  geniusFlashes: 0,
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
    const response = await apiFetch("/gamedata");
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
    get().connection?.close();
    set(resetAuthState);
    resetQueryDataState();
    await centralApiFetch("/auth/logout");
  },
  tap: (lotId: string) => {
    set((state) =>
      produce(state, (draftState) => {
        const lot = draftState.gameState.lots[lotId];
        if (lot) {
          lot.taps++;
        }
      })
    );
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

    const handleStateChange3 = (stateChange: GameStatePatch) => {
      unstable_batchedUpdates(() => {
        set((state) =>
          produce(state, (draftState) => {
            if (stateChange.gameState && stateChange.patchMask) {
              const partial = extractMessage(
                stateChange.gameState,
                stateChange.patchMask.paths
              );
              reflectionMergePartial(GameState, draftState.gameState, partial);
            } else if (stateChange.gameState) {
              draftState.gameState = stateChange.gameState;
            }
          })
        );
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
          handleStateChange3(msg.payload.stateChange);
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
        case "worldState":
          queryClient.setQueryData("worldState", msg.payload.worldState);
          break;
      }
    };
    const onStateChange = (connectionState: ConnectionState) => {
      unstable_batchedUpdates(() => {
        set((state) => {
          if (connectionState.error === "unauthorized") {
            resetQueryDataState();
            state = resetAuthState(state);
          }
          return { ...state, connectionState };
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
          },
        },
      })
    );
  },
  reschool: (mouseId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "reschoolMouse",
          reschoolMouse: {
            mouseId,
          },
        },
      })
    );
  },
  renameMouse: (mouseId, name) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "renameMouse",
          renameMouse: {
            mouseId,
            name,
          },
        },
      })
    );
  },
  openQuest: (questId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "openQuest",
          openQuest: {
            questId,
          },
        },
      })
    );
  },
  claimQuestReward: (questId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "claimQuestReward",
          claimQuestReward: {
            questId,
          },
        },
      })
    );
  },
  completeQuest: (questId) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "completeQuest",
          completeQuest: {
            questId,
          },
        },
      })
    );
  },
  reportActivity: () => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "reportActivity",
          reportActivity: {},
        },
      })
    );
  },
  saveMouseAppearance: (mouseId: string, appearance: MouseAppearance) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "saveMouseAppearance",
          saveMouseAppearance: {
            mouseId,
            appearance,
          },
        },
      })
    );
  },
  setAmbassadorMouse: (mouseId: string) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "setAmbassadorMouse",
          setAmbassadorMouse: {
            mouseId,
          },
        },
      })
    );
  },
  setPizzaPrice: (pizzaPrice: number) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "setPizzaPrice",
          setPizzaPrice: {
            pizzaPrice,
          },
        },
      })
    );
  },
  buyGeniusFlash: (id: number) => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "buyGeniusFlash",
          buyGeniusFlash: {
            id,
          },
        },
      })
    );
  },
}));

(window as any).useStore = useStore;
