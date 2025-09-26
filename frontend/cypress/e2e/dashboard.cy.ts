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

    
    cy.contains("Run Simulation").should("exist");


    
    cy.contains("View Map").should("exist");

    
  });

  it("displays recent simulations table with correct statuses", () => {
    cy.contains("Recent Simulations").should("exist");

    cy.get("table").within(() => {
      cy.contains("#1234").should("exist");
      cy.contains("Main St & 5th Ave").should("exist");
      cy.contains("Complete").should("exist");

      cy.contains("#1233").should("exist");
      cy.contains("Running").should("exist");

      cy.contains("#1232").should("exist");
      cy.contains("Failed").should("exist");

      cy.get("button").contains("View Details").should("have.length.at.least", 1);
    });
  });

  it("renders top intersections with progress bars", () => {
    cy.contains("Top Intersections").should("exist");

    cy.get(".intersection-item").should("have.length", 3);

    cy.get(".intersection-item").first().within(() => {
      cy.contains("Main St & 5th Ave").should("exist");
      cy.contains("15,000 vehicles").should("exist");
      cy.get(".progress-bar").should("have.css", "width");
    });

    cy.contains("Avg Daily Volume").should("exist");
    cy.contains("12,000 vehicles").should("exist");
  });
  
});

