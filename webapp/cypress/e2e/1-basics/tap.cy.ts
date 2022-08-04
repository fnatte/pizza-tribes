import { Building } from "../../../src/generated/building";
import { GameStatePatch } from "../../../src/generated/gamestate";

describe("tap", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          lots: {
            "1": {
              building: Building.SHOP,
            },
            "3": {
              building: Building.KITCHEN,
            },
          },
        },
        patchMask: {
          paths: ["lots.1", "lots.3"],
        },
      })
    );

    cy.visit("/");
  });

  afterEach(() => {
    cy.adminTestTeardown();
  });

  it("can tap 1 shop time", () => {
    // Tap
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="tap-section"] button').click();

    // Assert
    cy.get('[data-cy="tap-section"] button').should("be.disabled");
    cy.get('[data-cy="tap-section"]').contains("1 of 10");
    cy.get('[data-cy="tap-streak"]').should("have.attr", "aria-valuenow", 0);
    cy.get('[data-cy="tap-section"] button').should("not.be.disabled");
    cy.get('[data-cy="resource-bar-coins"]').should("have.text", 35);
  });

  it("can tap 10 shop times", () => {
    // Tap
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    for (let n = 0; n < 10; n++) {
      cy.get('[data-cy="tap-section"] button').click();
    }

    // Assert
    cy.get('[data-cy="tap-section"] button').should("be.disabled");
    cy.get('[data-cy="tap-section"]').contains("10 of 10");
    cy.get('[data-cy="tap-streak"]').should("have.attr", "aria-valuenow", 1);
    cy.get('[data-cy="resource-bar-coins"]').should("have.text", 350);
  });

  it("can tap 1 kitchen time", () => {
    // Tap
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot3"]').click();
    cy.get('[data-cy="tap-section"] button').click();

    // Assert
    cy.get('[data-cy="tap-section"] button').should("be.disabled");
    cy.get('[data-cy="tap-section"]').contains("1 of 10");
    cy.get('[data-cy="tap-streak"]').should("have.attr", "aria-valuenow", 0);
    cy.get('[data-cy="tap-section"] button').should("not.be.disabled");
    cy.get('[data-cy="resource-bar-pizzas"]').should("have.text", 80);
  });

  it("can tap 10 kitchen times", () => {
    // Tap
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot3"]').click();
    for (let n = 0; n < 10; n++) {
      cy.get('[data-cy="tap-section"] button').click();
    }

    // Assert
    cy.get('[data-cy="tap-section"] button').should("be.disabled");
    cy.get('[data-cy="tap-section"]').contains("10 of 10");
    cy.get('[data-cy="tap-streak"]').should("have.attr", "aria-valuenow", 1);
    cy.get('[data-cy="resource-bar-pizzas"]').should("have.text", 800);
  });
});
