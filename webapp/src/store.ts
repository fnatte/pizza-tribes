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
  Mouse,
  GameStatePatch_MousePatch,
  QuestState,
  GameStatePatch_QuestStatePatch,
  JsonPatchOp,
} from "./generated/gamestate";
import { GameData } from "./generated/game_data";
import { Report } from "./generated/report";
import { ResearchDiscovery } from "./generated/research";
import {
  ServerMessage,
  ServerMessage_GameStatePatch2,
  ServerMessage_Reports,
  ServerMessage_User,
} from "./generated/server_message";
import { Stats } from "./generated/stats";
import { generateId } from "./utils";
import { applyPatches, enablePatches, produce } from "immer";
import { queryClient } from "./queryClient";
import { apiFetch } from "./api";

enablePatches();

export type Lot = {
  building: Building;
  tappedAt: string;
  level: number;
  taps: number;
  streak: number;
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
  mice: Record<string, Mouse>;
  quests: Record<string, QuestState>;
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
  reschool: (mouseId: string) => void;
  renameMouse: (mouseId: string, name: string) => void;
  openQuest: (questId: string) => void;
  claimQuestReward: (questId: string) => void;
  completeVisitHelpPageQuest: () => void;
};

const resetQueryDataState = () => {
  queryClient.setQueriesData({}, () => undefined);
};

const mergeLots = (
  a: GameState["lots"],
  b: Record<string, GameStatePatch_LotPatch>
): GameState["lots"] => {
  const res = { ...a };

  Object.keys(b).forEach((lotId) => {
    const lot = res[lotId];
    if (b[lotId] !== undefined) {
      const { building, tappedAt, level, razed, taps, streak } = b[lotId];

      if (razed) {
        delete res[lotId];
        return;
      }

      if (lot === undefined) {
        res[lotId] = { building, tappedAt, level, taps, streak };
      } else {
        res[lotId] = {
          ...lot,
          building: building,
          tappedAt: tappedAt,
          level: level,
          taps: taps,
          streak: streak,
        };
      }
    }
  });

  return res;
};

function updateMice(
  mice: Record<string, Mouse>,
  updatedMice: { [key: string]: GameStatePatch_MousePatch }
): Record<string, Mouse> {
  const entries = Object.entries(mice)
    .map(([id, m]) => {
      const m2: Mouse = { ...m };
      const update = updatedMice[id];
      if (update) {
        if (update.name) {
          m2.name = update.name.value;
        }
        if (update.education) {
          m2.education = update.education.value;
        }
        if (update.isEducated) {
          m2.isEducated = update.isEducated.value;
        }
        if (update.isBeingEducated) {
          m2.isBeingEducated = update.isBeingEducated.value;
        }
      }
      return [id, m2] as const;
    })
    .filter(([id]) => updatedMice[id] !== null)
    .concat(
      Object.keys(updatedMice)
        .filter((id) => !mice[id])
        .map((id) => {
          const m: Mouse = {
            name: updatedMice[id].name?.value ?? "",
            isEducated: updatedMice[id].isEducated?.value ?? false,
            isBeingEducated: updatedMice[id].isBeingEducated?.value ?? false,
            education: updatedMice[id].education?.value ?? Education.CHEF,
          };
          return [id, m];
        })
    );

  return Object.fromEntries(entries);
}

function updateQuests(
  questStates: Record<string, QuestState>,
  questPatches: { [key: string]: GameStatePatch_QuestStatePatch }
): Record<string, QuestState> {
  const entries = Object.entries(questStates)
    .map(([id, q]) => {
      const q2: QuestState = { ...q };
      const patch = questPatches[id];
      if (patch) {
        if (patch.opened) {
          q2.opened = patch.opened.value;
        }
        if (patch.completed) {
          q2.completed = patch.completed.value;
        }
        if (patch.claimedReward) {
          q2.claimedReward = patch.claimedReward.value;
        }
      }
      return [id, q2] as const;
    })
    .filter(([id]) => questStates[id] !== null)
    .concat(
      Object.keys(questPatches)
        .filter((id) => !questStates[id])
        .map((id) => {
          const q: QuestState = {
            opened: questPatches[id].opened?.value ?? false,
            claimedReward: questPatches[id].claimedReward?.value ?? false,
            completed: questPatches[id].completed?.value ?? false,
          };
          return [id, q];
        })
    );

  return Object.fromEntries(entries);
}

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
  mice: {},
  quests: {},
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
    await apiFetch("/auth/logout");
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
            mice: updateMice(state.gameState.mice, stateChange.mice),
            quests: updateQuests(state.gameState.quests, stateChange.quests),
          },
        }));
      });
    };

    const isSupportedJsonPatchOp = (
      p: JsonPatchOp
    ): p is Omit<JsonPatchOp, "op"> & { op: "replace" | "remove" | "add" } => {
      return p.op === "replace" || p.op === "remove" || p.op === "add";
    };

    const handleStateChange2 = (stateChange: ServerMessage_GameStatePatch2) => {
      unstable_batchedUpdates(() => {
        set((state) => ({
          gameState:
            (
            applyPatches(
              state.gameState,
              stateChange.jsonPatch.filter(isSupportedJsonPatchOp).map((p) => {
                return {
                  path: p.path.split("/").filter((x) => x),
                  op: p.op,
                  value: JSON.parse(p.value),
                };
              })
            )),
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
        case "stateChange2":
          handleStateChange2(msg.payload.stateChange2);
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
  completeVisitHelpPageQuest: () => {
    get().connection?.send(
      ClientMessage.create({
        id: generateId(),
        type: {
          oneofKind: "completeVisitHelpPageQuest",
          completeVisitHelpPageQuest: {},
        },
      })
    );
  },
}));

(window as any).useStore = useStore;
