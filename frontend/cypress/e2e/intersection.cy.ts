/// <reference types="cypress" />

describe("Intersections Page", () => {
  const API_BASE_URL = "http://localhost:9090";

  beforeEach(() => {
    // Stub auth token in localStorage
    window.localStorage.setItem("authToken", "fake-jwt-token");
    cy.intercept("GET", `${API_BASE_URL}/intersections`, {
      statusCode: 200,
      body: {
        intersections: [
          {
            id: "123",
            name: "Test Intersection",
            traffic_density: "medium",
            details: { address: "Main Rd", city: "Pretoria", province: "Gauteng" },
            default_parameters: {
              optimisation_type: "default",
              simulation_parameters: {
                intersection_type: "trafficlight",
                green: 30,
                yellow: 3,
                red: 27,
                speed: 60,
                seed: 42,
              },
            },
          },
        ],
      },
    }).as("getIntersections");

    cy.visit("/intersections");
    cy.wait("@getIntersections");
  });

  it("should display fetched intersections", () => {
    cy.contains("Test Intersection").should("exist");
  });

  it("should open create intersection modal and submit form", () => {
    cy.contains("Create New Intersection").click();

    cy.get("input[name='name']").type("New Intersection");
    cy.get("input[name='details.address']").type("Corner of Lynnwood & Atterbury");
    cy.get("input[name='details.city']").clear().type("Pretoria");
    cy.get("input[name='details.province']").clear().type("Gauteng");

    cy.intercept("POST", `${API_BASE_URL}/intersections`, {
      statusCode: 201,
      body: { id: "456", name: "New Intersection" },
    }).as("createIntersection");

    cy.contains("Create Intersection").click();
    cy.wait("@createIntersection");

    cy.contains("New Intersection").should("exist");
  });

  it("should edit an intersection", () => {
    cy.contains("Test Intersection").parent().within(() => {
      cy.contains("Edit").click();
    });

    cy.get("input[name='name']").clear().type("Updated Intersection");

    cy.intercept("PATCH", `${API_BASE_URL}/intersections/123`, {
      statusCode: 200,
    }).as("updateIntersection");
    
  });
});

