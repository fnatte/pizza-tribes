import { AppearanceCategory } from "../../../src/generated/appearance";
import { GameState, GameStatePatch } from "../../../src/generated/gamestate";

describe("change appearance", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          mice: {
            "1": {
              name: "Mickey",
            },
            "2": {
              name: "Mus",
              appearance: {
                parts: {
                  [AppearanceCategory.HAT]: {
                    id: "redHat1",
                  },
                },
              },
            },
          },
        },
        patchMask: {
          paths: ["mice.1", "mice.2"],
        },
      })
    );
  });

  afterEach(() => {
    cy.adminTestTeardown();
  });

  it("can set ambassador", () => {
    cy.visit("/mouse/1");
    cy.get('[data-cy="make-ambassador-button"]').click();
    cy.get('[data-cy="make-ambassador-button"]').should("be.disabled");
  });

  it("can look at own ambassador", () => {
    cy.visit("/mouse/2");
    cy.get('[data-cy="make-ambassador-button"]').click();

    cy.adminGetGameState().as("gameState");
    cy.get<GameState>("@gameState").then((gameState) => {
      cy.visit(`/world/entry?x=${gameState.townX}&y=${gameState.townY}`);
      cy.get('[data-cy="ambassador-mouse"]').should("exist");

      const ambassador = gameState.mice[gameState.ambassadorMouseId];
      cy.wrap(ambassador)
        .its("appearance.parts")
        .should("have.nested.property", `${AppearanceCategory.HAT}.id`, "redHat1");
    });
  });
});
