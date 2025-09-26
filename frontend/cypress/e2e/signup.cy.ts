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

  it("registers successfully with valid input", () => {
    cy.intercept("POST", "http://localhost:9090/register", {
      statusCode: 200,
      body: {
        message: "Registration successful",
      },
    }).as("register");

    cy.get("input[name='username']").type("newuser");
    cy.get("input[name='email']").type("newuser@example.com");
    cy.get("input[name='password']").type("mypassword123");

    cy.contains("Register").click();

    cy.wait("@register");
    cy.contains("Registration successful! Redirecting to login...").should("exist");

    // Wait for redirect
    cy.url({ timeout: 5000 }).should("include", "/login");
  });

  it("shows error on failed registration", () => {
    cy.intercept("POST", "http://localhost:9090/register", {
      statusCode: 400,
      body: {
        message: "Email already exists",
      },
    }).as("registerFail");

    cy.get("input[name='username']").type("existinguser");
    cy.get("input[name='email']").type("existing@example.com");
    cy.get("input[name='password']").type("password123");

    cy.contains("Register").click();

    cy.wait("@registerFail");
    cy.contains("Email already exists").should("exist");
  });

  // it("updates traffic light based on inputs", () => {
  //   // cy.get(".traffic-light").within(() => {
  //   //   cy.get("div").eq(0).should("have.class", "bg-red-900/50");
  //   //   cy.get("div").eq(1).should("have.class", "bg-yellow-900/50");
  //   //   cy.get("div").eq(2).should("have.class", "bg-green-900/50");
  //   // });

  //   // cy.get("input[name='username']").type("user");
  //   // cy.get(".traffic-light").within(() => {
  //   //   cy.get("div").eq(0).should("have.class", "bg-red-600");
  //   // });

  //   // cy.get("input[name='email']").type("user@example.com");
  //   // cy.get(".traffic-light").within(() => {
  //   //   cy.get("div").eq(1).should("have.class", "bg-yellow-500");
  //   // });

  //   // cy.get("input[name='password']").type("secret123");
  //   // cy.get(".traffic-light").within(() => {
  //   //   cy.get("div").eq(2).should("have.class", "bg-green-500");
  //   // });
  // });

  it("navigates to login page when clicking Login here", () => {
    cy.contains("Login here").click();
    cy.url().should("include", "/login");
  });
});
