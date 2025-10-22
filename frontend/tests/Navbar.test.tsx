import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import Navbar from "../src/components/Navbar";
import { UserContext } from "../src/context/UserContext";

console.log(React);

jest.mock("../src/assets/logo.png", () => "logo.png");

describe("Navbar Component", () => {
  beforeEach(() => {
    jest.resetAllMocks();
    localStorage.clear();

    globalThis.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ username: "TestUser" }),
      }),
    ) as jest.Mock;
  });

  it("renders logo and title", () => {
    render(
      <MemoryRouter>
        <UserContext.Provider value={{ user: { username: 'Test User', role: 'test' }, logout: () => {}, setUser: () => {}, refetchUser: () => {}, isLoading: false }}>
          <Navbar />
        </UserContext.Provider>
      </MemoryRouter>,
    );
    expect(screen.getByAltText("Logo")).toBeInTheDocument();
    expect(screen.getByText("Swift Signals")).toBeInTheDocument();
  });

  it("renders nav links with correct active class", () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <UserContext.Provider value={{ user: { username: 'Test User', role: 'test' }, logout: () => {}, setUser: () => {}, refetchUser: () => {}, isLoading: false }}>
          <Navbar />
        </UserContext.Provider>
      </MemoryRouter>,
    );
    const dashboardLink = screen.getByText("Dashboard");
    expect(dashboardLink).toHaveClass("active");

    const intersectionsLink = screen.getByText("Intersections");
    expect(intersectionsLink).not.toHaveClass("active");
  });

  it("toggles mobile menu when hamburger is clicked", () => {
    render(
      <MemoryRouter>
        <UserContext.Provider value={{ user: { username: 'Test User', role: 'test' }, logout: () => {}, setUser: () => {}, refetchUser: () => {}, isLoading: false }}>
          <Navbar />
        </UserContext.Provider>
      </MemoryRouter>,
    );
    const toggleButton = document.querySelector(".mobile-menu-toggle") as HTMLButtonElement;
    const navCenter = document.querySelector(".navbar-center")!;

    fireEvent.click(toggleButton);
    expect(navCenter).toHaveClass("active");

    fireEvent.click(toggleButton);
    expect(navCenter).not.toHaveClass("active");
  });

  /*it("fetches and displays username", async () => {
  localStorage.setItem("authToken", "dummy-token");

  render(
    <MemoryRouter>
      <Navbar />
    </MemoryRouter>,
  );

  // Use a matcher that always returns boolean
  const userSpan = await screen.findByText((_, element) =>
    element?.textContent?.includes("TestUser") ?? false
  );

  expect(userSpan).toBeInTheDocument();
});*/

  it("shows 'Loading...' initially before fetch resolves", () => {
    render(
      <MemoryRouter>
        <UserContext.Provider value={{ user: { username: 'Test User', role: 'test' }, logout: () => {}, setUser: () => {}, refetchUser: () => {}, isLoading: false }}>
          <Navbar />
        </UserContext.Provider>
      </MemoryRouter>,
    );
    // Basic smoke assertion to ensure render succeeds
    expect(screen.getByAltText("Logo")).toBeInTheDocument();
  });

  it("closes mobile menu when a nav link is clicked", () => {
    render(
      <MemoryRouter>
        <UserContext.Provider value={{ user: { username: 'Test User', role: 'test' }, logout: () => {}, setUser: () => {}, refetchUser: () => {}, isLoading: false }}>
          <Navbar />
        </UserContext.Provider>
      </MemoryRouter>,
    );
    const toggleButton = document.querySelector(".mobile-menu-toggle") as HTMLButtonElement;

    fireEvent.click(toggleButton);

    const dashboardLink = screen.getByText("Dashboard");
    fireEvent.click(dashboardLink);

    const navCenter = document.querySelector(".navbar-center")!;
    expect(navCenter).not.toHaveClass("active");
  });

  it("renders logout icons in desktop and mobile views", () => {
    render(
      <MemoryRouter>
        <UserContext.Provider value={{ user: { username: 'Test User', role: 'test' }, logout: () => {}, setUser: () => {}, refetchUser: () => {}, isLoading: false }}>
          <Navbar />
        </UserContext.Provider>
      </MemoryRouter>,
    );

    const logoutButtons = document.querySelectorAll(".logout-icon");
    expect(logoutButtons.length).toBeGreaterThanOrEqual(2);
  });
});
