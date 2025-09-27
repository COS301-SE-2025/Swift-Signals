describe("Simulation Results Page", () => {
  beforeEach(() => {
    // Stub the API call
    cy.intercept("GET", "/intersections/*/simulate", {
      fixture: "simulationResults.json",
    }).as("getSimulationResults");

    cy.visit("/simulation-results/1"); // adjust route
    cy.wait("@getSimulationResults");
  });

  it("renders the page title and description", () => {
    cy.contains("Results for Corner of Albertus Street & Simon Vermooten Road").should("be.visible");
    cy.contains("Viewing detailed results for simulation").should("be.visible");
  });


});

