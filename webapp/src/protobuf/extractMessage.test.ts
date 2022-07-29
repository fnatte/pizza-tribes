import { GameState } from "../generated/gamestate";
import { ResearchDiscovery } from "../generated/research";
import { extractMessage } from "./extractMessage";

describe("extractMessage", () => {
  test("extracts single field", () => {
    expect(
      extractMessage(
        GameState.create({
          townX: 5,
          townY: 5,
        }),
        ["townX"]
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
        ["resources.coins"]
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
        ["resources.coins", "townX"]
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
        ["lots.5.taps"]
      )
    ).toEqual({
      lots: {
        5: {
          taps: 1,
        },
      },
    });
  });

  test("extracts undefined from nested map", () => {
    expect(
      extractMessage(
        GameState.create({
          lots: {},
        }),
        ["lots.5"]
      )
    ).toEqual({
      lots: {
        5: undefined,
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
        ["discoveries"]
      )
    ).toEqual({
      discoveries: [ResearchDiscovery.WEBSITE, ResearchDiscovery.MOBILE_APP],
    });
  });

  test("throws if trying to extract repeated field that is not in the last position of the path", () => {
    expect(() =>
      extractMessage(
        GameState.create({
          researchQueue: [{ completeAt: "123" }],
        }),
        ["researchQueue.completeAt"]
      )
    ).toThrowError();
  });
});
