describe("Simulation Results Page", () => {
  beforeEach(() => {
    // Stub the API call
    cy.intercept("GET", "/intersections/*/simulate", {
      fixture: "simulationResults.json",
    }).as("getSimulationResults");

    cy.visit("/simulation-results/1"); // adjust route
    cy.wait("@getSimulationResults");
  });


});

