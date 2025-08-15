describe("SignUp Page", () => {
  beforeEach(() => {
    cy.visit("/signup"); // Adjust if route is different
  });

  it("renders logo, title, and input fields", () => {
    cy.get('img[alt="Swift Signals Logo"]').should("be.visible");
    cy.contains("Welcome to Swift Signals").should("exist");
    cy.contains("Sign Up").should("exist");
    cy.get("input[name='username']").should("exist");
    cy.get("input[name='email']").should("exist");
    cy.get("input[name='password']").should("exist");
    cy.contains("Register").should("exist");
  });

  it("shows validation error on empty submit", () => {
    cy.contains("Register").click();
    cy.contains("Please fill in all fields.").should("exist");
  });

});

