describe("Login Page", () => {
  beforeEach(() => {
    cy.visit("/login"); // Assumes routing and dev server are ready
  });

  it("renders logo and form elements", () => {
    cy.get('img[alt="Swift Signals Logo"]').should("be.visible");
    cy.contains("Welcome to Swift Signals").should("exist");
    cy.contains("Login").should("exist");
    cy.get("input[name='username']").should("exist");
    cy.get("input[name='password']").should("exist");
    cy.contains("Log Me In").should("exist");
  });

  it("shows validation error on empty form submit", () => {
    cy.contains("Log Me In").click();
    cy.contains("Please fill in all fields.").should("exist");
  });

  
});

