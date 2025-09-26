/// <reference types="cypress" />

describe("Intersections Page", () => {
  const API_BASE_URL = "http://localhost:9090";

  beforeEach(() => {
    // Stub auth token in localStorage
    window.localStorage.setItem("authToken", "fake-jwt-token");
    cy.intercept("GET", `${API_BASE_URL}/intersections`, {
      statusCode: 200,
     
    }).as("getIntersections");
  });
});

