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

  it("submits form successfully with valid credentials", () => {
    // Stub backend response
    cy.intercept("POST", "http://localhost:9090/login", {
      statusCode: 200,
      body: {
        token: "mock-token",
        message: "Login successful",
      },
    }).as("loginRequest");

    cy.get("input[name='username']").type("testuser@example.com");
    cy.get("input[name='password']").type("securepassword");
    cy.contains("Log Me In").click();

    // Wait for mock request and confirm redirect
    cy.wait("@loginRequest");
    cy.url().should("include", "/dashboard");
  });

  it("shows login error on failed login", () => {
    cy.intercept("POST", "http://localhost:9090/login", {
      statusCode: 401,
      body: { message: "Invalid credentials" },
    }).as("loginFail");

    cy.get("input[name='username']").type("wrong@example.com");
    cy.get("input[name='password']").type("wrongpassword");
    cy.contains("Log Me In").click();

    cy.wait("@loginFail");
    cy.contains('Login failed. Server says: "Invalid credentials"').should("exist");
  });

  
});

