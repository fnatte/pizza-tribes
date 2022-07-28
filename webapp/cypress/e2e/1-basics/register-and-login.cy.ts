const username = "cypress-test-user";
const password = "secret";

describe("pizza-tribes", () => {
  beforeEach(() => {
    cy.visit("/");
    cy.adminDeleteUser(username);
  });

  it("register and login", () => {
    cy.contains("button", "Create Account").click();
    cy.location("pathname").should("eq", "/create-account");
    cy.get('input[name="username"').type(username);
    cy.get('input[name="password"').type(password);
    cy.get('input[name="confirm"').type(password);
    cy.get('button[type="submit"').click();
    cy.location("pathname").should("eq", "/login");

    cy.login(username, password)
  });

});
