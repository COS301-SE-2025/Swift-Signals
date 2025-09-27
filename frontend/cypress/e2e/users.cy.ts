describe("Users Page", () => {
    const API_BASE_URL = "http://localhost:9090";
    
  beforeEach(() => {
    // Stub token
    window.localStorage.setItem("authToken", "fake-admin-token");

    // Mock users API
    cy.intercept("GET", `${API_BASE_URL}/admin/users*`, [
      {
        id: "1",
        username: "Bob Johnson",
        email: "bob@example.com",
        is_admin: false,
        intersection_ids: [],
      },
      {
        id: "2",
        username: "Alice Smith",
        email: "alice@example.com",
        is_admin: false,
        intersection_ids: [],
      },
    ]).as("getUsers");

    cy.visit("/users");
    cy.wait("@getUsers");
  });

  it("should display users in the table", () => {
    cy.contains("User Management").should("exist");
    cy.contains("Bob Johnson").should("exist");
    cy.contains("Alice Smith").should("exist");
  });

    it("should open edit modal, update user, and close", () => {
        cy.contains("Alice Smith")
          .parent("tr")
          .within(() => {
            cy.get("button").contains(/edit/i).click();
        });

        cy.get("input#username").clear().type("Alice Updated");
        cy.get("input#email").clear().type("alice.updated@example.com");

        cy.intercept("PATCH", `${API_BASE_URL}/admin/users/2`, {
          statusCode: 200,
          body: {
            id: "2",
            username: "Alice Updated",
            email: "alice.updated@example.com",
            is_admin: false,
            intersection_ids: [],
          },
        }).as("updateUser");

        cy.contains("Save Changes").click();
        cy.wait("@updateUser");

        cy.contains("Alice Updated").should("exist");
  });

    it("should show error if update fails", () => {
        
    });
    
  // it("displays correct user information", () => {
  //   cy.get("table.usersTable tbody")
  //     .contains("td", "John Doe")
  //     .should("exist");
  //   cy.get("table.usersTable tbody")
  //     .contains("td", "Admin")
  //     .should("exist");
  // });

  // it("has working edit and delete buttons for each row", () => {
  //   cy.get("table.usersTable tbody tr").each(($row) => {
  //     cy.wrap($row).within(() => {
  //       cy.get("button[aria-label='Edit user']").should("exist").click();
  //       cy.get("button[aria-label='Delete user']").should("exist").click();
  //     });
  //   });
  // });

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
