import { Building } from "../../../src/generated/building";
import { GameStatePatch } from "../../../src/generated/gamestate";
import { ResearchDiscovery } from "../../../src/generated/research";

describe("research", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          lots: {
            "1": {
              building: Building.RESEARCH_INSTITUTE,
            },
          },
        },
        patchMask: {
          paths: ["lots.1"],
        },
      })
    );
  });

  afterEach(() => {
    cy.adminTestTeardown();
  });

  it("can research durum wheat", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: { geniusFlashes: 1 },
        patchMask: {
          paths: ["geniusFlashes"],
        },
      })
    );
    cy.gameVisit("/town/1");

    cy.get('[data-cy="ongoing-research-row"]').should("not.exist");

    cy.contains('[data-cy="research-area"]', "Demand").click();
    cy.contains('[data-cy="research-node"]', "Durum Wheat").click();

    cy.get('[data-cy="start-research-button"]').click();
    cy.get('[data-cy="ongoing-research-row"]').should("exist");

    cy.contains('[data-cy="research-area"]', "Demand").click();
    cy.contains('[data-cy="research-node"]', "Durum Wheat").click();
    cy.contains("Researching...").should("exist");

    cy.go(-1);

    cy.adminCompleteQueues();

    cy.get('[data-cy="ongoing-research-row"]').should("not.exist");
    cy.get('[data-cy="research-area"]').should("contain.text", "1 of");

    cy.contains('[data-cy="research-area"]', "Demand").click();
    cy.contains('[data-cy="research-node"]', "Durum Wheat").click();
    cy.contains("Already researched").should("exist");
  });

  it("can not research next node before previous are researched", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: { geniusFlashes: 1 },
        patchMask: {
          paths: ["geniusFlashes"],
        },
      })
    );
    cy.gameVisit("/town/1");

    cy.contains('[data-cy="research-area"]', "Demand").click();
    cy.contains('[data-cy="research-node"]', "Double Zero Flour").click();

    cy.get('[data-cy="start-research-button"]').should("be.disabled");
  });

  it("can research next node when previous is researched", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          geniusFlashes: 1,
          discoveries: [ResearchDiscovery.DURUM_WHEAT],
        },
        patchMask: {
          paths: ["geniusFlashes", "discoveries"],
        },
      })
    );
    cy.gameVisit("/town/1");

    cy.contains('[data-cy="research-area"]', "Demand").click();
    cy.contains('[data-cy="research-node"]', "Double Zero Flour").click();

    cy.get('[data-cy="start-research-button"]')
      .should("not.be.disabled")
      .click();
    cy.get('[data-cy="ongoing-research-row"]').should("exist");
  });

  it("can not research if has no genius flashes", () => {
    cy.gameVisit("/town/1");

    cy.contains('[data-cy="research-area"]', "Demand").click();
    cy.contains('[data-cy="research-node"]', "Durum Wheat").click();

    cy.get('[data-cy="start-research-button"]').should("be.disabled");
  });

  it("can buy a genius flash", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: { resources: { coins: 10_000, pizzas: 5_000 } },
        patchMask: {
          paths: ["resources.coins", "resources.pizzas"],
        },
      })
    );

    cy.gameVisit("/town/1");

    cy.get('[data-cy="get-more-button"]').click();
    cy.get('[data-cy="available-genius-flashes"]').should("have.text", "0");
    cy.get('[data-cy="buy-genius-flash-button"]')
      .should("not.be.disabled")
      .click();
    cy.get('[data-cy="available-genius-flashes"]').should("have.text", "1");
  });

  it("can not buy a genius flash if not enough resources", () => {
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: { resources: { coins: 100, pizzas: 100 } },
        patchMask: {
          paths: ["resources.coins", "resources.pizzas"],
        },
      })
    );

    cy.gameVisit("/town/1");

    cy.get('[data-cy="get-more-button"]').click();
    cy.get('[data-cy="available-genius-flashes"]').should("have.text", "0");
    cy.get('[data-cy="buy-genius-flash-button"]').should("be.disabled");
  });
});
