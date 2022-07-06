import Ajv, { JTDDataType } from "ajv/dist/jtd";
const ajv = new Ajv();

const gameStateSchema = {
  properties: {
    resources: {
      properties: {
        coins: { type: "int32" },
        pizzas: { type: "int32" },
      },
    },
    lots: {
      values: {
        properties: {
          building: { ref: "building" },
          tappedAt: { type: "int32" },
          level: { type: "int32" },
          taps: { type: "int32" },
          streak: { type: "int32" },
        },
      },
    },
    population: {
      properties: {
        uneducated: { type: "int32" },
        chefs: { type: "int32" },
        salesmice: { type: "int32" },
        guards: { type: "int32" },
        thieves: { type: "int32" },
        publicists: { type: "int32" },
      },
    },
    timestamp: { type: "int64" },
    trainingQueue: {
      elements: {
        properties: {
          completeAt: { type: "int64" },
          education: {
            ref: "education",
          },
          amount: { type: "int32" },
        },
      },
    },
    constructionQueue: {
      elements: {
        properties: {
          completeAt: { type: "int64" },
          lotId: { type: "string" },
          building: { ref: "building" },
          level: { type: "int32" },
          razing: { type: "bool" },
        },
      },
    },
    townX: { type: "int32" },
    townY: { type: "int32" },
    travelQueue: {
      elements: {
        properties: {
          arrivalAt: { type: "int64" },
          destinationX: { type: "int32" },
          destinationY: { type: "int32" },
          returning: { type: "bool" },
          thieves: { type: "int32" },
          coins: { type: "int64" },
        },
      },
    },
    discoveries: {
      elements: {
        ref: "research_discovery",
      },
    },
    researchQueue: {
      elements: {
        properties: {
          completeAt: { type: "int64" },
          discovery: { ref: "research_discovery" },
        },
      },
    },
    mice: {
      values: {
        properties: {
          name: { type: "string" },
          isEducated: { type: "bool" },
          isBeingEducated: { type: "bool" },
          education: { ref: "education" },
        },
      },
    },
    quests: {
      values: {
        properties: {
          opened: { type: "bool" },
          completed: { type: "bool" },
          claimedReward: { type: "bool" },
        },
      },
    },
  },
  definitions: {
    building: {
      enum: [
        "kitchen",
        "shop",
        "house",
        "school",
        "marketinghq",
        "research_institute",
        "town_centre",
      ],
    },
    education: {
      enum: ["chef", "salesmouse", "guard", "thief", "publicist"],
    },
    research_discovery: {
      enum: [
        "website",
        "digital_ordering_system",
        "mobile_app",
        "masonry_oven",
        "gas_oven",
        "hybrid_oven",
        "durum_wheat",
        "double_zero_flour",
        "san_marzano_tomatoes",
        "ocimum_basilicum",
        "extra_virgin",
      ],
    },
  },
} as const;

export type Education = JTDDataType<typeof gameStateSchema.definitions.education>;
export type Building = JTDDataType<typeof gameStateSchema.definitions.building>;
export type GameState = JTDDataType<typeof gameStateSchema>;
export type Construction = GameState['constructionQueue'][0]
export type Lot = GameState['lots']['']

export const validateGameState = ajv.compile<GameState>(gameStateSchema);
