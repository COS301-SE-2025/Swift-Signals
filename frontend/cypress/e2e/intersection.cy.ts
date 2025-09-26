/// <reference types="cypress" />

describe("Intersections Page", () => {
  const API_BASE_URL = "http://localhost:9090";

  beforeEach(() => {
    // Stub auth token in localStorage
    window.localStorage.setItem("authToken", "fake-jwt-token");
  });
});

