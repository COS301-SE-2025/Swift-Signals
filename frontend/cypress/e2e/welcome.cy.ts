describe("Welcome Page", () => {
  beforeEach(() => {
    cy.visit("http://localhost:5173/"); // Or use env baseUrl
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
});

