describe("Users Page", () => {
  beforeEach(() => {
    cy.visit("/users"); // Adjust path as necessary
  });

  it("renders a table with 9 user rows", () => {
    cy.get("table.usersTable").should("exist");
    cy.get("table.usersTable tbody tr").should("have.length", 9);
  });

  it("displays correct user information", () => {
    cy.get("table.usersTable tbody")
      .contains("td", "John Doe")
      .should("exist");
    cy.get("table.usersTable tbody")
      .contains("td", "Admin")
      .should("exist");
  });
  
});

