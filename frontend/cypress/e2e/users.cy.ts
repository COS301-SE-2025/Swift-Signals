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

  it("has working edit and delete buttons for each row", () => {
    cy.get("table.usersTable tbody tr").each(($row) => {
      cy.wrap($row).within(() => {
        cy.get("button[aria-label='Edit user']").should("exist").click();
        cy.get("button[aria-label='Delete user']").should("exist").click();
      });
    });
  });

  // it("paginates to the next page", () => {
  //   cy.get("tbody tr").first().should("contain", "John Doe"); // Page 1
  //   cy.get("button[aria-label='Next page']").click();
  //   cy.get("tbody tr").first().should("not.contain", "John Doe"); // Page 2
  // });

  // it("paginates back to the previous page", () => {
  //   cy.get("button[aria-label='Next page']").click();
  //   cy.get("button[aria-label='Previous page']").click();
  //   cy.get("tbody tr").first().should("contain", "John Doe");
  // });

  // it("displays ellipsis for long pagination", () => {
  //   cy.get(".usersPaging").should("contain", "...");
  // });

  // it("navigates directly to a specific page", () => {
  //   cy.get(".usersPaging button").contains("2").click();
  //   cy.get("tbody tr").first().within(() => {
  //     cy.get("td").eq(1).should("exist"); // Name cell exists
  //   });
  // });
});
