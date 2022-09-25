/// <reference types="cypress" />

import { GameState, GameStatePatch } from "../../src/generated/gamestate";

const defaultUsername = "cypress_test_user";
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
      adminBatchCreateUser(
        users: { username: string; password: string }[]
      ): Chainable<void>;
      adminBatchDeleteUser(usernames: string[]): Chainable<void>;
      adminTestSetup(username?: string): Chainable<void>;
      adminTestTeardown(username?: string): Chainable<void>;

      login(username?: string, password?: string): Chainable<void>;
      expand(): Chainable<Element>;
      gameVisit(url: string): Chainable<AUTWindow>;
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
  cy.request<string[]>(
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
        .request("GET", `http://localhost:8081/users/${userId}/gameState`)
        .then((r2) => {
          expect(r2.isOkStatusCode).to.be.true;
        })
        .its("body")
        .then<GameState>((x) => GameState.fromJson(x));
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

Cypress.Commands.add("adminBatchCreateUser", (users) => {
  cy.request<{ users: { status: number }[] }>(
    "POST",
    "http://localhost:8081/users/batch",
    {
      users,
    }
  ).then((x) => {
    expect(x.status).to.eq(207);
    x.body.users.forEach((user) => {
      expect(user.status).to.be.within(200, 299);
    });
  });
});

Cypress.Commands.add("adminBatchDeleteUser", (usernames) => {
  cy.request<{ users: { status: number }[] }>(
    "DELETE",
    "http://localhost:8081/users/batch",
    {
      usernames,
    }
  ).then((x) => {
    expect(x.status).to.eq(207);
    x.body.users.forEach((user) => {
      expect(user.status).to.be.satisfy(
        (status: number) => (status >= 200 && status < 300) || status === 404
      );
    });
  });
});

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
      cy.wrap(resp.body.gameId).as("gameId");

      window.localStorage.setItem(
        "CapacitorStorage.accessToken",
        resp.body.accessToken
      );
    });
});

Cypress.Commands.add("adminTestTeardown", (username = defaultUsername) => {
  cy.adminBatchDeleteUser([username]);
});

Cypress.Commands.add(
  "login",
  (username = defaultUsername, password = defaultPassword) => {
    cy.visit("/");
    cy.location("pathname").should("eq", "/login");
    cy.get('input[name="username"').type(username);
    cy.get('input[name="password"').type(password);
    cy.get('button[type="submit"').click();
    cy.location("pathname").should("eq", "/games");
  }
);

Cypress.Commands.add("expand", { prevSubject: true }, (subject) => {
  cy.wrap(subject).then(($el) => {
    if ($el.parents('[aria-expanded="true"]').length === 0) {
      cy.wrap($el).click();
    }
  });
});

Cypress.Commands.add("gameVisit", (url: string) => {
  cy.get<string>("@gameId").then((gameId) => {
    cy.visit(`/game/${gameId}${url}`);
  });
});

export {};
