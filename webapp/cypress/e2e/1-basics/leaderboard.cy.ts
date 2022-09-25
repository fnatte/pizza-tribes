import { GameStatePatch } from "../../../src/generated/gamestate";

describe("leaderboard", () => {
  beforeEach(() => {
    cy.adminTestSetup();

    cy.adminBatchDeleteUser(
      [...Array(30).keys()].map((i) => `cypress_other_test_user_${i + 1}`)
    );
    cy.adminBatchCreateUser(
      [...Array(30).keys()].map((i) => ({
        username: `cypress_other_test_user_${i + 1}`,
        password: "test",
      }))
    );

    for (let i = 0; i < 30; i++) {
      const username = `cypress_other_test_user_${i + 1}`;
      cy.adminPatchGameState(
        GameStatePatch.create({
          gameState: {
            resources: {
              coins: i * 1000,
            },
          },
          patchMask: {
            paths: ["resources.coins"],
          },
        }),
        username
      );
    }

    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          resources: {
            coins: 4000,
          },
        },
        patchMask: {
          paths: ["resources.coins"],
        },
      })
    );
  });

  afterEach(() => {
    cy.adminTestTeardown();
    cy.adminBatchDeleteUser(
      [...Array(30).keys()].map((i) => `cypress_other_test_user_${i + 1}`)
    );
  });

  it("can see leaderboard", () => {
    cy.gameVisit("/leaderboard");
    cy.get('[data-cy="leaderboard-row"]')
      .contains("cypress_test_user")
      .closest('[data-cy="leaderboard-row"]')
      .within(() => {
        cy.root().should("contain.text", "4,000");
        cy.get('td').eq(0).invoke('text').then(parseFloat).should("be.gte", 27);
      });
  });
});
