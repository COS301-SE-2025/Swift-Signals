describe("Welcome Page", () => {
  beforeEach(() => {
    cy.visit("/"); // set in cypress.config.cjs
  });

  it("should display the heading and logo", () => {
    cy.contains("Welcome to Swift Signals!").should("exist");
    cy.get('img[alt="Logo"]').should("be.visible");
  });

  it("should navigate to Login page on Login button click", () => {
    cy.contains("Login").click();
    cy.url().should("include", "/login");
  });

  it("should navigate to Register page on Register button click", () => {
    cy.contains("Register").click();
    cy.url().should("include", "/signup");
  });

  describe("Carousel interaction", () => {
    it("should render the carousel and show the first item", () => {
      cy.get(".carousel-container").should("be.visible");

      // First slide (title: "Overview") should be visible
      cy.get(".carousel-item").first().within(() => {
        cy.get(".carousel-item-title").should("contain.text", "Overview");
      });

    });

  });
  
});

