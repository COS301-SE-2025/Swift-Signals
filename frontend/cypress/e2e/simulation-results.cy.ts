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

  it("displays the KPI cards with values", () => {
    const kpis = [
      "Avg Speed",
      "Avg Delay",
      "Avg Flow Rate",
      "Emissions",
      "# Phases",
      "Cycle Time",
      "Safety Severity",
    ];

    kpis.forEach((label) => {
      cy.contains(label).should("exist");
    });

    // Verify at least some numeric values exist
    cy.get(".kpi-card").each(($card) => {
      cy.wrap($card).invoke("text").should("match", /\d/);
    });
  });

  it("shows the chart section", () => {
    cy.contains("Simulation Results vs Optimized Results").should("be.visible");
    cy.get("canvas, svg").should("have.length.at.least", 4);
  });


});

