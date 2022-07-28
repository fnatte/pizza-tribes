/// <reference types="cypress" />

const defaultUsername = "cypress-test-user";
const defaultPassword = "secret";

declare global {
  namespace Cypress {
    interface Chainable {
      adminCompleteQueues(username?: string): Chainable<void>;
      adminIncrCoins(amount: number, username?: string): Chainable<void>;
      adminDeleteUser(username?: string): Chainable<void>;
      adminCreateUser(username?: string, password?: string): Chainable<void>;
      adminTestSetup(username?: string): Chainable<void>;
      login(username?: string, password?: string): Chainable<void>;
    }
  }
}

Cypress.Commands.add("adminDeleteUser", (username = defaultUsername) => {
  cy.request("GET", `http://localhost:8081/users?username=${username}`).then(
    (r) => {
      const users = r.body;
      if (users.length === 0) {
        return;
      }

      const userId = users[0];
      cy.request("DELETE", `http://localhost:8081/users/${userId}`).then(
        (r2) => {
          expect(r2.isOkStatusCode).to.be.true;
        }
      );
    }
  );
});

Cypress.Commands.add("adminCompleteQueues", (username = defaultUsername) => {
  cy.request("GET", `http://localhost:8081/users?username=${username}`).then(
    (r) => {
      const users = r.body;
      if (users.length === 0) {
        return;
      }

      const userId = users[0];
      cy.request(
        "POST",
        `http://localhost:8081/users/${userId}/completeQueues`
      ).then((r2) => {
        expect(r2.isOkStatusCode).to.be.true;
      });
    }
  );
});

Cypress.Commands.add(
  "adminCreateUser",
  (username = defaultUsername, password = defaultPassword) => {
    cy.request("POST", "http://localhost:8081/users", {
      username,
      password,
    });
  }
);

Cypress.Commands.add("adminIncrCoins", (amount, username = defaultUsername) => {
  cy.request("GET", `http://localhost:8081/users?username=${username}`).then(
    (r) => {
      const users = r.body;
      if (users.length === 0) {
        return;
      }

      const userId = users[0];
      cy.request("POST", `http://localhost:8081/users/${userId}/incrCoins`, {
        amount,
      }).then((r2) => {
        expect(r2.isOkStatusCode).to.be.true;
      });
    }
  );
});

Cypress.Commands.add("adminTestSetup", (username = defaultUsername) => {
  cy.request("POST", "http://localhost:8081/test/setup", {
    username,
  })
    .as("testSetup")
    .then((resp) => {
      window.localStorage.setItem(
        "CapacitorStorage.accessToken",
        resp.body.accessToken
      );
    });
});

Cypress.Commands.add(
  "login",
  (username = defaultUsername, password = defaultPassword) => {
    cy.visit("/");
    cy.location("pathname").should("eq", "/login");
    cy.get('input[name="username"').type(username);
    cy.get('input[name="password"').type(password);
    cy.get('button[type="submit"').click();
    cy.location("pathname").should("eq", "/");
  }
);

export {};
