import { Building } from "../../../src/generated/building";
import { GameStatePatch } from "../../../src/generated/gamestate";

describe("education", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminPatchGameState(
      GameStatePatch.create({
        gameState: {
          resources: {
            coins: 100_000,
          },
          lots: {
            "1": {
              building: Building.SCHOOL,
            },
            "3": {
              building: Building.HOUSE,
            },
          },
          mice: {
            "1": {},
            "2": {},
            "3": {},
            "4": {},
            "5": {},
          },
        },
        patchMask: {
          paths: ["resources.coins", "lots.1", "lots.3", "mice"],
        },
      })
    );

    cy.visit("/");
  });

  it("can train 1 chef", () => {
    // Train
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="school-education-title"]')
      .contains("Chef")
      .parents('[data-cy="school-education"]')
      .find('[data-cy="school-education-submit-button"]')
      .click();

    // Assert
    cy.get('[data-cy="training-queue-table-row"]').should(($tr) => {
      expect($tr.find('[data-cy="training-queue-table-amount"]').text()).to.eq(
        "1"
      );
      expect($tr.find('[data-cy="training-queue-table-title"]').text()).to.eq(
        "Chef"
      );
    });
  });

  it("can train 3 salesmice", () => {
    // Train
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="school-education-title"]')
      .contains("Salesmouse")
      .parents('[data-cy="school-education"]')
      .find('[data-cy="school-education-amount-input"]')
      .clear()
      .type("3");
    cy.get('[data-cy="school-education-title"]')
      .contains("Salesmouse")
      .parents('[data-cy="school-education"]')
      .find('[data-cy="school-education-submit-button"]')
      .click();

    // Assert
    cy.get('[data-cy="training-queue-table-row"]').should(($tr) => {
      expect($tr.find('[data-cy="training-queue-table-amount"]').text()).to.eq(
        "3"
      );
      expect($tr.find('[data-cy="training-queue-table-title"]').text()).to.eq(
        "Salesmouse"
      );
    });
  });

  it("can't train more chefs than available uneducated", () => {
    // Train
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="school-education-title"]')
      .contains("Chef")
      .parents('[data-cy="school-education"]')
      .find('[data-cy="school-education-amount-input"]')
      .clear()
      .type("6");
    cy.get('[data-cy="school-education-title"]')
      .contains("Chef")
      .parents('[data-cy="school-education"]')
      .find('[data-cy="school-education-submit-button"]')
      .click();

    // Assert
    cy.get('[data-cy="school-education-title"]')
      .contains("Chef")
      .parents('[data-cy="school-education"]')
      .find('[data-cy="error"]')
      .should("exist");
  });
});
