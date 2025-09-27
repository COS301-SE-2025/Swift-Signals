describe("TrafficSimulation E2E Tests", () => {
  const INTERSECTION_ID = "test-intersection";

  beforeEach(() => {
    // Set up auth token and intercepts before visiting
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { fixture: "simulation.json" }
    ).as("getSim");

    cy.visit("/simulation-results", {
      onBeforeLoad: (win) => {
        win.localStorage.setItem("authToken", "fake-token");
      }
    });
  });

  it("shows loading state", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { delay: 1000, fixture: "simulation.json" }
    ).as("getSimDelayed");

    // Navigate to the simulation (replace with your actual trigger)
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.contains("Loading simulation data...").should("be.visible");
  });

  it("shows error state if API fails", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { statusCode: 500 }
    ).as("getSimError");

    // Navigate to the simulation
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.contains("Error: Failed to fetch simulation data").should("be.visible");
    cy.contains("Retry").should("be.visible").click();
  });

  it("shows no data state", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { body: { output: null } }
    ).as("getSimEmpty");

    // Navigate to the simulation
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.contains("No simulation data available").should("be.visible");
  });

  it("renders simulation and UI controls", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { fixture: "simulation.json" }
    ).as("getSim");

    // Navigate to the simulation
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.wait("@getSim");

    // Root container exists
    cy.get(".traffic-simulation-root").should("exist");

    // Canvas should be rendered
    cy.get("canvas").should("exist");

    // SimulationUI controls
    cy.contains("Play").should("exist");
    cy.contains("Restart").should("exist");
  });

  it("handles missing auth token", () => {
    // Clear localStorage
    cy.clearLocalStorage();

    // Navigate to the simulation
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.contains("Error: Authentication token not found").should("be.visible");
  });

  it("handles 401 authentication error", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { statusCode: 401 }
    ).as("getSimUnauth");

    // Navigate to the simulation
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.contains("Error: Authentication failed. Please log in again.").should("be.visible");
  });

  it("handles 404 not found error", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/simulate`, 
      { statusCode: 404 }
    ).as("getSimNotFound");

    // Navigate to the simulation
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .viewBtn`).click();

    cy.contains("Error: Simulation data not found for this intersection.").should("be.visible");
  });

  it("works with optimize endpoint", () => {
    cy.intercept("GET", `http://localhost:9090/intersections/${INTERSECTION_ID}/optimise`, 
      { fixture: "simulation.json" }
    ).as("getOptimize");

    // Navigate to optimization (adjust selector as needed)
    cy.get(`[data-intersection-id="${INTERSECTION_ID}"] .optimizeBtn`).click();

    cy.wait("@getOptimize");
    cy.get(".traffic-simulation-root").should("exist");
  });
});
