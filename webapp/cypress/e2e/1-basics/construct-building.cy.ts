describe("constructing buildings", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.visit("/");
  });

  it("construct level 1 kitchen at lot 1", () => {
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot1"]').click();
    cy.get('[data-cy="construct-building-title"]').contains("Kitchen")
      .parents('[data-cy="construct-building"]')
      .find('[data-cy="construct-building-button"]')
      .click();

    cy.get('[data-cy="construction-queue-toggle-button"]').click();
    cy.get('[data-cy="construction-queue-table"] tr').should(
      "contain.text",
      "Kitchen"
    );

    cy.adminCompleteQueues();

    cy.get('[data-cy="lot1"] title').should("have.text", "Kitchen");
    cy.get('[data-cy="lot1"] [data-cy="level-badge"]').should("have.text", "1");
  });

  it("construct level 1 shop at lot 3", () => {
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot3"]').click();
    cy.get('[data-cy="construct-building-title"]').contains("Shop")
      .parents('[data-cy="construct-building"]')
      .find('[data-cy="construct-building-button"]')
      .click();

    cy.get('[data-cy="construction-queue-toggle-button"]').click();
    cy.get('[data-cy="construction-queue-table"] tr').should(
      "contain.text",
      "Shop"
    );

    cy.adminCompleteQueues();

    cy.get('[data-cy="lot3"] title').should("have.text", "Shop");
    cy.get('[data-cy="lot3"] [data-cy="level-badge"]').should("have.text", "1");
  });

  it("construct level 2 house at lot 4", () => {
    cy.adminIncrCoins(1000)
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot4"]').click();
    cy.get('[data-cy="construct-building-title"]').contains("House")
      .parents('[data-cy="construct-building"]')
      .find('[data-cy="construct-building-button"]')
      .click();
    cy.adminCompleteQueues();
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
    cy.get('[data-cy="lot4"]').click();
    cy.get('[data-cy="town-lot"]').should("contain.text", "House")
    cy.get('[data-cy="upgrade-building-button"]').click();
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();

    cy.get('[data-cy="construction-queue-toggle-button"]').click();
    cy.get('[data-cy="construction-queue-table"] tr').should(
      "contain.text",
      "House"
    );

    cy.adminCompleteQueues();

    cy.get('[data-cy="lot4"] title').should("have.text", "House");
    cy.get('[data-cy="lot4"] [data-cy="level-badge"]').should("have.text", "2");
  });

});
