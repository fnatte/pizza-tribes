describe("constructing buildings", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.visit("/");
  });

  afterEach(() => {
    cy.adminTestTeardown();
  });

  it("construct level 1 kitchen at lot 1", () => {
    cy.get('[data-cy="main-nav"] a[href$="/town"]').click();
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

    cy.get('[data-cy="lot1"]').should("have.attr", 'title', "Kitchen");
    cy.get('[data-cy="lot1"] [data-cy="level-badge"]').should("have.text", "1");
  });

  it("construct level 1 shop at lot 3", () => {
    cy.get('[data-cy="main-nav"] a[href$="/town"]').click();
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

    cy.get('[data-cy="lot3"]').should("have.attr", "title", "Shop");
    cy.get('[data-cy="lot3"] [data-cy="level-badge"]').should("have.text", "1");
  });

  it("construct level 2 house at lot 5", () => {
    cy.adminIncrCoins(1000)
    cy.get('[data-cy="main-nav"] a[href$="/town"]').click();
    cy.get('[data-cy="lot5"]').click();
    cy.get('[data-cy="construct-building-title"]').contains("House")
      .parents('[data-cy="construct-building"]')
      .find('[data-cy="construct-building-button"]')
      .click();
    cy.adminCompleteQueues();
    cy.get('[data-cy="main-nav"] a[href$="/town"]').click();
    cy.get('[data-cy="lot5"]').click();
    cy.get('[data-cy="town-lot"]').should("contain.text", "House")
    cy.get('[data-cy="upgrade-building-button"]').click();
    cy.get('[data-cy="main-nav"] a[href$="/town"]').click();

    cy.get('[data-cy="construction-queue-toggle-button"]').click();
    cy.get('[data-cy="construction-queue-table"] tr').should(
      "contain.text",
      "House"
    );

    cy.adminCompleteQueues();

    cy.get('[data-cy="lot5"]').should("have.attr", "title", "House");
    cy.get('[data-cy="lot5"] [data-cy="level-badge"]').should("have.text", "2");
  });

});
