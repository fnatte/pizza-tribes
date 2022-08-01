import { AppearanceCategory } from "../../../src/generated/appearance";
import { GameStatePatch } from "../../../src/generated/gamestate";

describe("change appearance", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          mice: { "1": {} },
        },
        patchMask: {
          paths: ["mice.1"],
        },
      })
    );

    cy.visit("/mouse/1/appearance");
  });

  it("can set outfit", () => {
    cy.get('[data-cy="gallery-section-title"]').contains("Outfits").expand();
    cy.get('[data-cy="gallery-section-item"]').first().click();
    cy.get('[data-cy="gallery-section-title"]').contains("Hats").expand();
    cy.get('[data-cy="gallery-section-item"]').first().click();

    cy.get('[data-cy="mouse-editor-save-button"]').click();
    cy.location("pathname").should("eq", "/mouse/1");

    cy.adminGetGameState().as("gameState")
    cy.get("@gameState")
      .its("mice.1.appearance.parts")
      .should("have.nested.property", `${AppearanceCategory.OUTFIT}.id`, "guardOutfit1")
    cy.get("@gameState")
      .its("mice.1.appearance.parts")
      .and("have.nested.property", `${AppearanceCategory.HAT}.id`, "chefHat1");
  });
});
