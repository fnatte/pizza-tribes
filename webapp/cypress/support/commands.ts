/// <reference types="cypress" />

import { GameState, GameStatePatch } from "../../src/generated/gamestate";

const defaultUsername = "cypress-test-user";
const defaultPassword = "secret";

declare global {
  namespace Cypress {
    interface Chainable {
      adminGetGameState(username?: string): Chainable<GameState | null>;
      adminCompleteQueues(username?: string): Chainable<void>;
      adminIncrCoins(amount: number, username?: string): Chainable<void>;
      adminPatchGameState(
        req: GameStatePatch,
        username?: string
      ): Chainable<void>;
      adminDeleteUser(username?: string): Chainable<void>;
      adminCreateUser(username?: string, password?: string): Chainable<void>;
      adminTestSetup(username?: string): Chainable<void>;

      login(username?: string, password?: string): Chainable<void>;
      expand(): Chainable<Element>;
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

Cypress.Commands.add("adminGetGameState", (username = defaultUsername) => {
  return cy
    .request<string[]>(
      "GET",
      `http://localhost:8081/users?username=${username}`
    )
    .its("body")
    .then((users) => {
      if (users.length === 0) {
        return cy.wrap(null);
      }

      const userId = users[0];
      return cy
        .request(
          "GET",
          `http://localhost:8081/users/${userId}/gameState`
        )
        .then((r2) => {
          expect(r2.isOkStatusCode).to.be.true;
        })
        .its("body")
        .then<GameState>(x => GameState.fromJson(x))
    });
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

Cypress.Commands.add(
  "adminPatchGameState",
  (req, username = defaultUsername) => {
    const body = GameStatePatch.toJson(req) as object;

    cy.request("GET", `http://localhost:8081/users?username=${username}`).then(
      (r) => {
        const users = r.body;
        if (users.length === 0) {
          return;
        }

        const userId = users[0];
        cy.request(
          "PATCH",
          `http://localhost:8081/users/${userId}/gameState`,
          body
        ).then((r2) => {
          expect(r2.isOkStatusCode).to.be.true;
        });
      }
    );
  }
);

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

Cypress.Commands.add("expand", { prevSubject: true }, (subject) => {
  cy.wrap(subject).then(($el) => {
      if ($el.parents('[aria-expanded="true"]').length === 0) {
        cy.wrap($el).click();
      }
    });
});

export {};
