// tests/Navbar.test.tsx
import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import Navbar from "../src/components/Navbar";

console.log(React)

// Mock react-router-dom's useLocation
jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useLocation: jest.fn(() => ({
    pathname: "/dashboard",
  })),
}));

// Mock the logo import to just return a string
jest.mock("../src/assets/logo.png", () => "logo.png");

describe("Navbar Component", () => {
  beforeEach(() => {
    jest.resetAllMocks();
    localStorage.clear();

    // Mock fetch for user profile using globalThis
    globalThis.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ username: "TestUser" }),
      })
    ) as jest.Mock;
  });

  it("renders logo and site name", () => {
    render(<Navbar />);
    expect(screen.getByAltText("Logo")).toBeInTheDocument();
    expect(screen.getByText("Swift Signals")).toBeInTheDocument();
  });

  it("renders nav links with correct active class", () => {
    render(<Navbar />);
    const dashboardLink = screen.getByText("Dashboard");
    expect(dashboardLink).toHaveClass("active");
    const intersectionsLink = screen.getByText("Intersections");
    expect(intersectionsLink).not.toHaveClass("active");
  });

  it("toggles mobile menu when hamburger is clicked", () => {
    render(<Navbar />);
    const toggleButton = screen.getByRole("button");
    fireEvent.click(toggleButton);
    const navCenter = document.querySelector(".navbar-center");
    expect(navCenter).toHaveClass("active");
    fireEvent.click(toggleButton);
    expect(navCenter).not.toHaveClass("active");
  });

  it("fetches and displays username", async () => {
    render(<Navbar />);
    await waitFor(() => {
      expect(screen.getAllByText("TestUser")[0]).toBeInTheDocument();
    });
  });

  it("shows 'Loading...' initially before fetch resolves", () => {
    render(<Navbar />);
    expect(screen.getAllByText("Loading...")[0]).toBeInTheDocument();
  });

  it("calls toggleMobileMenu when mobile nav link is clicked", () => {
    render(<Navbar />);
    const toggleButton = screen.getByRole("button");
    fireEvent.click(toggleButton);

    const dashboardLink = screen.getByText("Dashboard");
    fireEvent.click(dashboardLink);

    const navCenter = document.querySelector(".navbar-center");
    expect(navCenter).not.toHaveClass("active");
  });

  it("renders logout icons", () => {
    render(<Navbar />);
    const logoutIcons = screen.getAllByRole("link", { name: "" });
    expect(logoutIcons.length).toBeGreaterThanOrEqual(2);
  });
});
