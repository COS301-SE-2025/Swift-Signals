import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import Navbar from "../src/components/Navbar";

console.log(React)

jest.mock("../src/assets/logo.png", () => "logo.png");

describe("Navbar Component", () => {
  beforeEach(() => {
    jest.resetAllMocks();
    localStorage.clear();

    globalThis.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ username: "TestUser" }),
      })
    ) as jest.Mock;
  });

  it("renders logo and site name", () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Navbar />
      </MemoryRouter>
    );
    expect(screen.getByAltText("Logo")).toBeInTheDocument();
    expect(screen.getByText("Swift Signals")).toBeInTheDocument();
  });

  it("renders nav links with correct active class", () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Navbar />
      </MemoryRouter>
    );
    const dashboardLink = screen.getByText("Dashboard");
    expect(dashboardLink).toHaveClass("active");

    const intersectionsLink = screen.getByText("Intersections");
    expect(intersectionsLink).not.toHaveClass("active");
  });

  it("toggles mobile menu when hamburger is clicked", () => {
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );
    const toggleButton = screen.getByRole("button");
    const navCenter = document.querySelector(".navbar-center")!;

    fireEvent.click(toggleButton);
    expect(navCenter).toHaveClass("active");

    fireEvent.click(toggleButton);
    expect(navCenter).not.toHaveClass("active");
  });

  it("fetches and displays username", async () => {
    localStorage.setItem("authToken", "dummy-token");

    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getAllByText("TestUser")[0]).toBeInTheDocument();
    });
  });

  it("shows 'Loading...' initially before fetch resolves", () => {
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );
    expect(screen.getAllByText("Loading...")[0]).toBeInTheDocument();
  });

  it("closes mobile menu when a nav link is clicked", () => {
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );

    const toggleButton = screen.getByRole("button");
    fireEvent.click(toggleButton);

    const dashboardLink = screen.getByText("Dashboard");
    fireEvent.click(dashboardLink);

    const navCenter = document.querySelector(".navbar-center")!;
    expect(navCenter).not.toHaveClass("active");
  });

  it("renders logout icons in desktop and mobile views", () => {
    render(
      <MemoryRouter>
        <Navbar />
      </MemoryRouter>
    );

    const logoutIcons = screen.getAllByRole("link", { name: "" });
    expect(logoutIcons.length).toBeGreaterThanOrEqual(2);
  });
});
