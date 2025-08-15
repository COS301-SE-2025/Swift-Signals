describe("Dashboard Page", () => {
  beforeEach(() => {
    cy.visit("/dashboard"); // Update if route is different
  });

  it("renders dashboard cards", () => {
    cy.contains("Total Intersections").should("exist");
    cy.contains("24").should("exist");

    cy.contains("Active Simulations").should("exist");
    cy.contains("8").should("exist");

    cy.contains("Optimization Runs").should("exist");
    cy.contains("156").should("exist");
  });

  it("shows quick action buttons", () => {
    cy.contains("New Intersection").should("exist");
    cy.contains("Run Simulation").should("exist");
    cy.contains("View Map").should("exist");
  });
  
});

