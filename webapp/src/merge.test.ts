import { GameState } from "./generated/gamestate";
import { FieldMask } from "./generated/google/protobuf/field_mask";
import { ResearchDiscovery } from "./generated/research";
import { extractMessage } from "./merge";

describe("extractMessage", () => {
  test("extracts single field", () => {
    expect(
      extractMessage(
        GameState.create({
          townX: 5,
          townY: 5,
        }),
        FieldMask.create({
          paths: ["townX"],
        })
      )
    ).toEqual({
      townX: 5,
    });
  });

  test("extracts field nested in submessage", () => {
    expect(
      extractMessage(
        GameState.create({
          townX: 5,
          townY: 5,
          resources: {
            coins: 10,
          },
        }),
        FieldMask.create({
          paths: ["resources.coins"],
        })
      )
    ).toEqual({
      resources: {
        coins: 10,
      },
    });
  });

  test("extracts multiple paths", () => {
    expect(
      extractMessage(
        GameState.create({
          townY: 5,
          townX: 5,
          resources: {
            coins: 10,
            pizzas: 11,
          },
        }),
        FieldMask.create({
          paths: ["resources.coins", "townX"],
        })
      )
    ).toEqual({
      townX: 5,
      resources: {
        coins: 10,
      },
    });
  });

  test("extracts field nested in map", () => {
    expect(
      extractMessage(
        GameState.create({
          townY: 5,
          townX: 5,
          resources: {
            coins: 10,
            pizzas: 11,
          },
          lots: {
            5: {
              taps: 1,
              level: 5,
            },
          },
        }),
        FieldMask.create({
          paths: ["lots.5.taps"],
        })
      )
    ).toEqual({
      lots: {
        5: {
          taps: 1,
        },
      },
    });
  });

  test("extracts repeated field if it is the last position of the path", () => {
    expect(
      extractMessage(
        GameState.create({
          discoveries: [
            ResearchDiscovery.WEBSITE,
            ResearchDiscovery.MOBILE_APP,
          ],
        }),
        FieldMask.create({
          paths: ["discoveries"],
        })
      )
    ).toEqual({
      discoveries: [ResearchDiscovery.WEBSITE, ResearchDiscovery.MOBILE_APP],
    });
  });

  test("throws if trying to extract repeated field that is not in the last position of the path", () => {
    expect(
      extractMessage(
        GameState.create({
          researchQueue: [
            { completeAt: "123" }
          ]
        }),
        FieldMask.create({
          paths: ["researchQueue.0.completeAt"],
        })
      )
    ).toThrowError();
  });
});
