describe("ComparisonView Page", () => {
  const ORIGINAL_ID = "123";

  beforeEach(() => {
    cy.window().then((win) => {
      win.localStorage.setItem("authToken", "fake-token");
    });
  });

  it("renders both panels with labels", () => {
    cy.intercept("GET", `https://swiftsignals.seebranhome.co.za/intersections/${ORIGINAL_ID}/simulate`, {
      fixture: "simulation.json",
    }).as("getSim");

    cy.intercept("GET", `https://swiftsignals.seebranhome.co.za/intersections/${ORIGINAL_ID}/optimise`, {
      fixture: "optimized.json",
    }).as("getOpt");

    cy.visit(`/comparison/${ORIGINAL_ID}`);

    cy.wait("@getSim");

    cy.contains("Original Simulation").should("be.visible");
    cy.contains("Optimized Simulation").should("be.visible");

    cy.get(".comparison-view").should("have.length", 2);
  });

  it("can exit via Exit button and Escape key", () => {
    cy.visit(`/comparison/${ORIGINAL_ID}`);

    cy.contains("Exit").should("be.visible").click();
    // Depending on router setup, assert navigation or window close fallback
    cy.url().should("not.include", `/comparison/${ORIGINAL_ID}`);

    // Reload and test Escape key
    cy.visit(`/comparison/${ORIGINAL_ID}`);
    cy.get("body").type("{esc}");
    cy.url().should("not.include", `/comparison/${ORIGINAL_ID}`);
  });

  it("toggles fullscreen for left panel", () => {
    cy.visit(`/comparison/${ORIGINAL_ID}`);

    cy.contains("Fullscreen").first().click();
    cy.contains("Exit Fullscreen").should("be.visible");
    cy.contains("Exit Fullscreen").click();
    cy.contains("Fullscreen").should("be.visible");
  });

  it("shows loading state while checking optimized data", () => {
    cy.intercept("GET", `https://swiftsignals.seebranhome.co.za/intersections/${ORIGINAL_ID}/optimise`, {
      delay: 1000,
      fixture: "optimized.json",
    }).as("getOpt");

    cy.visit(`/comparison/${ORIGINAL_ID}`);

    cy.contains("Checking Optimization Status").should("be.visible");
    cy.wait("@getOpt");
  });

  it("shows 'No Optimization Available' when no optimized data", () => {
    cy.intercept("GET", `https://swiftsignals.seebranhome.co.za/intersections/${ORIGINAL_ID}/optimise`, {
      body: { output: { vehicles: [] } },
    }).as("getOpt");

    cy.visit(`/comparison/${ORIGINAL_ID}`);

    cy.contains("No Optimization Available").should("be.visible");
    cy.contains("Check for Optimization").should("be.visible").click();

    cy.wait("@getOpt");
  });

  it("renders optimized simulation when data exists", () => {
    cy.intercept("GET", `https://swiftsignals.seebranhome.co.za/intersections/${ORIGINAL_ID}/optimise`, {
      fixture: "optimized.json",
    }).as("getOpt");

    cy.visit(`/comparison/${ORIGINAL_ID}`);
    cy.wait("@getOpt");

    cy.contains("Optimized Simulation").should("be.visible");
    cy.contains("Fullscreen").last().click();
    cy.contains("Exit Fullscreen").should("be.visible");
  });

  it("displays error if optimization check fails", () => {
    cy.intercept("GET", `https://swiftsignals.seebranhome.co.za/intersections/${ORIGINAL_ID}/optimise`, {
      statusCode: 500,
    }).as("getOpt");

    cy.visit(`/comparison/${ORIGINAL_ID}`);
    cy.wait("@getOpt");

    cy.contains("No Optimization Available").should("be.visible");
    cy.contains("Error").should("be.visible");
  });
});

