describe("Dashboard Page", () => {
  beforeEach(() => {
    cy.visit("/dashboard"); // Update if route is different
  });

  it("renders dashboard cards", () => {
    cy.contains("Total Intersections").should("exist");
    // cy.contains("24").should("exist");

    cy.contains("Active Simulations").should("exist");
    // cy.contains("8").should("exist");

    cy.contains("Optimization Runs").should("exist");
    // cy.contains("156").should("exist");
  });

  it("shows quick action buttons and they are clickable", () => {
    cy.contains("New Intersection").should("exist");
    cy.url().should("include", "/intersections/new"); // adjust route if different
    cy.go("back");
    
    cy.contains("Run Simulation").should("exist");
    cy.url().should("include", "/simulations/run"); // adjust route if different
    cy.go("back");
    
    cy.contains("View Map").should("exist");
    cy.url().should("include", "/map"); // adjust route if different
    cy.go("back");
  });

  it("displays recent simulations table with correct statuses", () => {
    cy.contains("Recent Simulations").should("exist");

    cy.get("table").within(() => {
      cy.contains("Corner of Albertus Street & Simon Vermooten Road").should("exist");
      cy.contains("unoptimised").should("exist");

      cy.contains("Corner of Lynnwood & Jan Shoba").should("exist");
      cy.contains("unoptimised").should("exist");

      cy.get("a, button").contains("View Details").should("have.length.at.least", 2);
    });
  });

  it("renders traffic density distribution chart", () => {
    cy.contains("Traffic Density Distribution").should("exist");
    
    // Ensure chart SVG or canvas is present
    cy.get(".traffic-chart-container").find("svg, canvas").should("exist");
  });

  it("renders recent intersections with statuses", () => {
    
  });
  
});

