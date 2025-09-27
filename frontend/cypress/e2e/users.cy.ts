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
        cy.contains("Bob Johnson")
          .parent("tr")
          .within(() => {
            cy.get("button").contains(/edit/i).click();
        });

        cy.get("input#username").clear().type("Broken Bob");

        cy.intercept("PATCH", `${API_BASE_URL}/admin/users/1`, {
          statusCode: 400,
          body: { message: "Invalid data" },
        }).as("updateFail");

        cy.contains("Save Changes").click();
        cy.wait("@updateFail");

        cy.contains("Invalid data").should("exist");
    });

    it("should delete a user after confirmation", () => {
        cy.on("window:confirm", () => true); // accept confirm dialog

        cy.intercept("DELETE", `${API_BASE_URL}/admin/users/1`, {
          statusCode: 204,
        }).as("deleteUser");

        cy.contains("Bob Johnson")
          .parent("tr")
          .within(() => {
            cy.get("button").contains(/delete/i).click();
        });

        cy.wait("@deleteUser");

        cy.contains("Bob Johnson").should("not.exist");
        
    });

    it("should handle pagination", () => {
        // Mock 10 users
        cy.intercept("GET", `${API_BASE_URL}/admin/users*`, Array.from({ length: 10 }).map((_, i) => ({
          id: String(i + 1),
          username: `User${i + 1}`,
          email: `user${i + 1}@example.com`,
          is_admin: false,
          intersection_ids: [],
        }))).as("getManyUsers");
    
        cy.visit("/users");
        cy.wait("@getManyUsers");

        cy.contains("User1").should("exist");
        cy.contains("Next").click();
        cy.contains("User9").should("exist");
      });


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
