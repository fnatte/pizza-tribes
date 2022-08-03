import { GameStatePatch } from "../../../src/generated/gamestate";

describe("leaderboard", () => {
  beforeEach(() => {
    cy.adminTestSetup();

    cy.adminBatchDeleteUser(
      [...Array(30).keys()].map((i) => `cypress_test_user_${i + 1}`)
    );
    cy.adminBatchCreateUser(
      [...Array(30).keys()].map((i) => ({
        username: `cypress_test_user_${i + 1}`,
        password: "test",
      }))
    );

    for (let i = 0; i < 30; i++) {
      const username = `cypress_test_user_${i + 1}`;
      cy.adminPatchGameState(
        GameStatePatch.create({
          gameState: {
            resources: {
              coins: i * 1000,
            },
          },
          patchMask: {
            paths: ["lots.1", "mice"],
          },
        }),
        username
      );
    }
  });

  it("can see leaderboard", () => {
    cy.visit("/leaderboard");
  });
});
