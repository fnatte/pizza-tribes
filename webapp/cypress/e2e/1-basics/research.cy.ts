import { Building } from "../../../src/generated/building";
import { GameStatePatch } from "../../../src/generated/gamestate";

describe("research", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          resources: {
            coins: 1_000_000,
          },
          lots: {
            "1": {
              building: Building.RESEARCH_INSTITUTE,
            },
          },
        },
        patchMask: {
          paths: ["lots.1", "resources.coins"],
        },
      })
    );
    cy.visit("/town/1");
  });

  it("can research durum wheat", () => {
    cy.get('[data-cy="ongoing-research-row"]').should("not.exist");

    cy.get('[data-cy="research-track"]')
      .contains("Pizza Craft")
      .parents('[data-cy="research-track"]')
      .find('[data-cy="research-track-expand-toggle"]')
      .click();

    cy.get('[data-cy="research-node-expand-toggle"]:first').click();
    cy.get('[data-cy="research-node-start-research-button"]').click();
    cy.get('[data-cy="ongoing-research-row"]').should("exist");
    cy.get('[data-cy="research-node-being-researched"]').should("exist");

    cy.adminCompleteQueues();

    cy.get('[data-cy="ongoing-research-row"]').should("not.exist");
    cy.get('[data-cy="research-node-already-researched"]').should("exist");

    cy.get('[data-cy="research-track"]')
      .invoke("text")
      .should("match", /durum wheat has already been researched/i);

    cy.get('[data-cy="research-track"]').should("contain.text", "1 of");
  });

  it("can not research last node before previous are researched", () => {
    cy.get('[data-cy="research-track"]')
      .contains("Pizza Craft")
      .parents('[data-cy="research-track"]')
      .find('[data-cy="research-track-expand-toggle"]')
      .click();

    cy.get('[data-cy="research-node-expand-toggle"]').last().click();
    cy.get('[data-cy="research-node-parent-not-discovered"]').should('exist');
  });
});
