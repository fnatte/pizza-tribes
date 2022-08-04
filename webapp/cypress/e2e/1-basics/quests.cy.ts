import { Building } from "../../../src/generated/building";
import { Education } from "../../../src/generated/education";
import { GameStatePatch } from "../../../src/generated/gamestate";
import { MiceBuilder } from "../../support/helpers";

describe("quests", () => {
  beforeEach(() => {
    cy.adminTestSetup();

    cy.visit("/quests");
  });

  afterEach(() => {
    cy.adminTestTeardown();
  });

  it("can see first quest", () => {
    cy.get('[data-cy="quest-item-title"]').should(
      "contain.text",
      "Bake and sell"
    );
    cy.get('[data-cy="quest-item-reward-coins"]').should("contain.text", "550");
    cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");
  });

  it("can complete first quest", () => {
    cy.get('[data-cy="quest-item-title"]').should(
      "contain.text",
      "Bake and sell"
    );

    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          lots: {
            "6": {
              building: Building.SHOP,
            },
            "7": {
              building: Building.KITCHEN,
            },
          },
        },
        patchMask: {
          paths: ["lots.6", "lots.7"],
        },
      })
    );

    cy.get('[data-cy="quest-item-claim-reward-button"]').click();

    // Assert
    cy.get('[data-cy="resource-bar-coins"]').should("contain.text", "550");
    cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");
    cy.get('[data-cy="quest-item"]').should("contain.text", "claimed");

    // Next quest should be available
    cy.get('[data-cy="quest-item-title"]').should("contain.text", "Workforce");
  });

  it("can complete visit help page quest", () => {
    cy.get('[data-cy="quest-item-title"]').should(
      "contain.text",
      "Bake and sell"
    );

    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          quests: {
            "5": {
              opened: true,
              completed: true,
              claimedReward: true,
            },
          },
        },
        patchMask: {
          paths: ["quests"],
        },
      })
    );

    // The help quest should be available after quest 5 has been completed
    cy.get('[data-cy="quest-item-title"]').should("contain.text", "Knowledge");
    cy.get('[data-cy="quest-item-title"]').contains("Knowledge").click();
    cy.get('[data-cy="quest-item-reward-coins"]').should("contain.text", "750");

    // Open menu and navigate to help page
    cy.get('[data-cy="menu-expand-button"]').click();
    cy.get('[data-cy="main-nav"] a[href="/help"]').click();
    cy.get('[data-cy="game-help"]').should("be.visible");

    // Go back to quests and claim the reward
    cy.get('[data-cy="main-nav"] a[href="/quests"]').click();
    cy.get('[data-cy="quest-item-claim-reward-button"]').click();

    // Assert
    cy.get('[data-cy="resource-bar-coins"]').should("contain.text", "750");
    cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");
    cy.get('[data-cy="quest-item"]').should("contain.text", "claimed");
  });

  it("can complete stats quest", () => {
    cy.get('[data-cy="quest-item-title"]').should(
      "contain.text",
      "Bake and sell"
    );

    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          quests: {
            "7": {
              opened: true,
              completed: true,
              claimedReward: true,
            },
          },
        },
        patchMask: {
          paths: ["quests"],
        },
      })
    );

    cy.get('[data-cy="quest-item-title"]').should("contain.text", "Statistics");
    cy.get('[data-cy="quest-item-title"]').contains("Statistics").expand();
    cy.get('[data-cy="stats-quest-form"]').should("exist");

    // Goto stats page to find the answer
    cy.get('[data-cy="main-nav"] [href="/stats"]').click();
    cy.get('[data-cy="stats-row"]')
      .contains("Pizzas produced")
      .parents('[data-cy="stats-row"]')
      .find('[data-cy="stats-value"]')
      .invoke("text")
      .as("answer");

    // Go back to quests page and fill in answer
    cy.get('[data-cy="main-nav"] [href="/quests"]').click();
    cy.get('[data-cy="quest-item-title"]').contains("Statistics").expand();
    cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");

    // Fill in wrong answer
    cy.get('[data-cy="stats-quest-form"]').within(() => {
      cy.get("input").type("wrong answer");
      cy.get('[type="submit"]').click();
      cy.root().should("contain.text", "Wrong");
      cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");
    });

    // Fill in correct answer
    cy.get<string>("@answer").then((answer) => {
      cy.get('[data-cy="stats-quest-form"]').within(() => {
        cy.get("input").type(answer);
        cy.get('[type="submit"]').click();
      });
    });

    // Claim reward
    cy.get('[data-cy="quest-item-claim-reward-button"]').click();
    cy.get('[data-cy="resource-bar-coins"]').should("contain.text", "1,000");
    cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");
    cy.get('[data-cy="quest-item"]').should("contain.text", "claimed");
  });

  it("can complete 8 employees quest", () => {
    cy.get('[data-cy="quest-item-title"]').should(
      "contain.text",
      "Bake and sell"
    );

    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          quests: {
            "7": {
              opened: true,
              completed: true,
              claimedReward: true,
            },
          },
        },
        patchMask: {
          paths: ["quests"],
        },
      })
    );

    cy.get('[data-cy="quest-item-title"]').should("contain.text", "Work, work");
    cy.get('[data-cy="quest-item-title"]').contains("Work, work").expand();

    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          lots: {
            "6": {
              building: Building.SHOP,
              level: 4,
            },
            "7": {
              building: Building.KITCHEN,
              level: 4,
            },
          },
          mice: new MiceBuilder()
            .add(Education.CHEF, 3)
            .add(Education.SALESMOUSE, 5)
            .build(),
        },
        patchMask: {
          paths: ["lots.6", "lots.7", "mice"],
        },
      })
    );

    // Claim reward
    cy.get('[data-cy="quest-item-claim-reward-button"]').click();
    cy.get('[data-cy="quest-item-claim-reward-button"]').should("not.exist");
    cy.get('[data-cy="quest-item"]').should("contain.text", "claimed");
  });
});
