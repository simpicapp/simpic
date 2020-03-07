/// <reference types="cypress" />

describe("Login test", () => {
  beforeEach(() => {
    cy.server();

    cy.route({
      url: "/api/users/me",
      status: 401,
      response: "",
    });

    cy.route({
      url: "/api/timeline?offset=0",
      response: [],
    });
  });

  it("Shows login form when the button is clicked", () => {
    cy.visit("/");
    cy.contains("Login").click();
    cy.get("input#username").should("be.visible");
    cy.get("input#password").should("be.visible");
  });

  it("Shows alert when logging in fails", () => {
    cy.route({
      method: "POST",
      url: "/api/login",
      status: 401,
      response: {error: "Insufficiently amazing"},
    });

    cy.visit("/");
    cy.contains("Login").click();
    cy.get("input#username").type("tester");
    cy.get("input#password").type("testing123");
    cy.contains("input[type=submit]", "Login").click();
    cy.contains("Insufficiently amazing").should("have.class", "alert");
  });

  it("Shows alert when logging in succeeds", () => {
    cy.route({
      method: "POST",
      url: "/api/login",
      status: 200,
      response: {},
    });

    cy.visit("/");
    cy.contains("Login").click();
    cy.get("input#username").type("tester");
    cy.get("input#password").type("testing123");
    cy.contains("input[type=submit]", "Login").click();
    cy.contains("You are now logged in").should("be.visible");
  });

  it("Refreshes user info when logging in succeeds", () => {
    cy.route({
      method: "POST",
      url: "/api/login",
      status: 200,
      response: {},
    });

    cy.visit("/");
    cy.contains("Login").click();
    cy.get("input#username").type("tester");
    cy.get("input#password").type("testing123");

    cy.route({
      method: "GET",
      url: "/api/users/me",
      status: 200,
      response: {username: "cypress", is_admin: true},
    });

    cy.contains("input[type=submit]", "Login").click();
    cy.contains("Logged in as cypress").should("be.visible");
    cy.contains("Logout").should("be.visible");
  });
});
