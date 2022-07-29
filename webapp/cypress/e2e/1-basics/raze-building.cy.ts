import { Building } from "../../../src/generated/building";
import { GameStatePatch } from "../../../src/generated/gamestate";

describe("raze building", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          lots: {
            "1": {
              building: Building.SHOP,
            },
          },
        },
        patchMask: {
          paths: ["lots.1"],
        },
      })
    );

    cy.visit("/");
  });

  it("can raze level 1 shop", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: { resources: { coins: 10_000 } },
        patchMask: {
          paths: ["resources.coins"],
        },
      })
    );

    // Raze
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="raze-building-button"]').click();

    cy.adminCompleteQueues();

    // Assert
    cy.get('[data-cy="cancel-raze-building-button"]').should("exist");
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"] svg').should("not.exist");
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="construct-buildings"]').should("exist");
  });

  it("cannot raze if insufficient coins", () => {
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="raze-building-button"]').should("be.disabled");
  });

  it("can cancel raze level 1 shop", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: { resources: { coins: 10_000 } },
        patchMask: {
          paths: ["resources.coins"],
        },
      })
    );

    // Raze
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="raze-building-button"]').click();
    cy.get('[data-cy="cancel-raze-building-button"]').click();

    cy.adminCompleteQueues();

    // Assert
    cy.get('[data-cy="cancel-raze-building-button"]').should("not.exist");
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"] svg').should("exist");
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="raze-building-button"]').should("exist");
  });
});
