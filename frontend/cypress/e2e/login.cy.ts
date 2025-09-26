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

  it("opens forgot password modal and submits reset request", () => {
    // Open modal
    cy.contains("Forgot Password?").click();
    cy.contains("Reset Password").should("exist");

    cy.get("input[name='resetEmail']").type("reset@example.com");

    cy.intercept("POST", "http://localhost:9090/reset-password", {
      statusCode: 200,
      body: {
        message: "Password reset instructions sent to your email.",
      },
    }).as("resetRequest");

    cy.contains("Send Reset Link").click();
    cy.wait("@resetRequest");
    cy.contains("Password reset instructions sent to your email.").should("exist");
  });

  it("shows error for invalid reset email", () => {
    cy.contains("Forgot Password?").click();

    cy.get("input[name='resetEmail']").type("fail@example.com");
    cy.intercept("POST", "http://localhost:9090/reset-password", {
      statusCode: 400,
      body: { message: "Reset failed. Email not found." },
    }).as("resetFail");

    cy.contains("Send Reset Link").click();
    cy.wait("@resetFail");
    cy.contains("Reset failed. Email not found.").should("exist");
  });

  it("updates traffic light based on form input", () => {
    // cy.get(".traffic-light").within(() => {
    //   // Initially all off
    //   cy.get("div").eq(0).should("have.class", "bg-red-900/50");
    //   cy.get("div").eq(1).should("have.class", "bg-yellow-900/50");
    //   cy.get("div").eq(2).should("have.class", "bg-green-900/50");
    // });

    cy.get("input[name='username']").type("user");
    cy.get(".traffic-light").within(() => {
      cy.get("div").eq(0).should("have.class", "bg-red-600");
      cy.get("div").eq(1).should("have.class", "bg-yellow-500");
    });

    cy.get("input[name='password']").type("pass");
    cy.get(".traffic-light").within(() => {
      cy.get("div").eq(2).should("have.class", "bg-green-500");
    });
  });

  it("navigates to signup on Register Here button click", () => {
    cy.contains("Register Here").click();
    cy.url().should("include", "/signup");
  });
});
