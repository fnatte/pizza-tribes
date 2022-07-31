import { Building } from "../../../src/generated/building";
import { Education } from "../../../src/generated/education";
import { GameState, GameStatePatch, Mouse } from "../../../src/generated/gamestate";

describe("travel", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminDeleteUser("target");
    cy.adminCreateUser("target");
    cy.adminGetGameState("target").as("targetGameState");
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          lots: {
            "1": {
              building: Building.HOUSE,
            },
          },
          mice: {
            "1": {
              isEducated: true,
              name: "1",
              education: Education.THIEF,
            },
            "2": {
              isEducated: true,
              name: "2",
              education: Education.THIEF,
            },
            "3": {
              isEducated: true,
              name: "3",
              education: Education.THIEF,
            },
          },
        },
        patchMask: {
          paths: ["lots.1", "mice"],
        },
      })
    );

    cy.get<GameState>("@targetGameState").then((gameState) => {
      cy.visit(`/world/entry?x=${gameState.townX}&y=${gameState.townY}`);
    });
  });

  it("can send 1 thief", () => {
    // Send
    cy.get('[data-cy="thieves-to-send-input"]').clear().type("1");
    cy.get('[data-cy="send-thieves-button"]').click();
    cy.get('[data-cy="travel-queue-toggle-button"]').click();
    cy.get('[data-cy="travel-queue-row"]').should(
      "contain.text",
      "1 travelling thief"
    );

    // Returning
    cy.adminCompleteQueues();

    // Check report
    cy.get('[data-cy="menu-expand-button"]').click();
    cy.get('[data-cy="main-nav"] a[href="/reports"]').click();
    cy.get('[data-cy="report-row"]').should("have.length", 1);
    cy.get('[data-cy="report-link"]').should("have.length", 1).click();
    cy.get('[data-cy="report-title"]')
      .should("have.text", "Thief report")
      .click();
    cy.get('[data-cy="show-report"]')
      .should("contain.text", "0 coins")
      .invoke("text")
      .should("match", /no thieves were caught/i);

    // Check travel queue
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="travel-queue-toggle-button"]').click();
    cy.get('[data-cy="travel-queue-row"]').should(
      "contain.text",
      "1 returning thief"
    );

    // Back home
    cy.adminCompleteQueues();

    // Check travel queue
    cy.get('[data-cy="travel-queue-row"]').should("not.exist");
  });

  it("can send 3 thieves", () => {
    // Send
    cy.get('[data-cy="thieves-to-send-input"]').clear().type("3");
    cy.get('[data-cy="send-thieves-button"]').click();
    cy.get('[data-cy="travel-queue-toggle-button"]').click();
    cy.get('[data-cy="travel-queue-row"]').should(
      "contain.text",
      "3 travelling thieves"
    );
    cy.adminCompleteQueues();

    // Returning
    cy.get('[data-cy="travel-queue-row"]').should(
      "contain.text",
      "3 returning thieves"
    );
    cy.adminCompleteQueues();

    // Back home
    cy.get('[data-cy="travel-queue-row"]').should("not.exist");
  });

  it("can not send 4 thieves", () => {
    cy.get('[data-cy="thieves-to-send-input"]').clear().type("4");

    cy.get('[data-cy="error"]').should("not.exist");
    cy.get('[data-cy="send-thieves-button"]').click();
    cy.get('[data-cy="error"]').should("exist");
  });

  it("can steal coins", () => {
    // Give target 10 000 coins
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          resources: {
            coins: 10_000,
          },
        },
        patchMask: {
          paths: ["resources.coins"],
        },
      }),
      "target"
    );

    // Send
    cy.get('[data-cy="thieves-to-send-input"]').clear().type("3");
    cy.get('[data-cy="send-thieves-button"]').click();
    cy.adminCompleteQueues();

    // Check report
    cy.get('[data-cy="menu-expand-button"]').click();
    cy.get('[data-cy="main-nav"] a[href="/reports"]').click();
    cy.get('[data-cy="report-row"]').should("have.length", 1);
    cy.get('[data-cy="report-link"]').should("have.length", 1).click();
    cy.get('[data-cy="report-title"]')
      .should("have.text", "Thief report")
      .click();
    cy.get('[data-cy="show-report"]')
      .should("contain.text", "9,000 coins")
      .invoke("text")
      .should("match", /no thieves were caught/i);

    // Returning
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="travel-queue-toggle-button"]').click();
    cy.get('[data-cy="travel-queue-row"]').should(
      "contain.text",
      "3 returning thieves"
    );
    cy.adminCompleteQueues();

    // Back home
    cy.get('[data-cy="resource-bar-coins"]').should("have.text", "9,000");

    // Target should have lost 9000 coins
    cy.adminGetGameState("target").its('resources.coins').should('eq', 1_000);
  });

  it("should lose thieves if there are guards", () => {
    // Give target 1000 guards
    const mice: Record<string, Partial<Mouse>> = {};
    for (let n = 0; n < 1000; n++) {
      mice[`m${n}`] = {
        isEducated: true,
        education: Education.GUARD,
      }
    }
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          mice
        },
        patchMask: {
          paths: ["mice"],
        },
      }),
      "target"
    );

    // Send
    cy.get('[data-cy="thieves-to-send-input"]').clear().type("3");
    cy.get('[data-cy="send-thieves-button"]').click();
    cy.adminCompleteQueues();

    // Check report
    cy.get('[data-cy="menu-expand-button"]').click();
    cy.get('[data-cy="main-nav"] a[href="/reports"]').click();
    cy.get('[data-cy="report-row"]').should("have.length", 1);
    cy.get('[data-cy="report-link"]').should("have.length", 1).click();
    cy.get('[data-cy="report-title"]')
      .should("have.text", "Thief report")
      .click();
    cy.get('[data-cy="show-report"]')
      .should("contain.text", "3 thieves got caught");

    // Returning
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="travel-queue-toggle-button"]').should('not.exist');
    cy.get('[data-cy="population-table-toggle-button"]').click();
    cy.get('[data-cy="population-table-uneducated-count"]').should("have.text", "3");
  });
});

