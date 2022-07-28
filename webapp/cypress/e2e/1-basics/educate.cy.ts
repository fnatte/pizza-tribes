describe("education", () => {
  beforeEach(() => {
    cy.adminTestSetup();
    cy.adminIncrCoins(100_000);
    cy.visit("/");

    cy.adminCompleteQueues();
  });

  it("train 1 chef", () => {
    cy.get('[data-cy="main-nav"] a[href="/town"]').click();
  });

});
