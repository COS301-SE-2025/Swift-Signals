// tests/ThemeToggle.test.tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import ThemeToggle from "../src/components/ThemeToggle";

console.log(React)

describe("ThemeToggle Component", () => {
  const localStorageMock = (() => {
    let store: Record<string, string> = {};
    return {
      getItem: (key: string) => store[key] || null,
      setItem: (key: string, value: string) => {
        store[key] = value;
      },
      clear: () => {
        store = {};
      },
      removeItem: (key: string) => {
        delete store[key];
      },
    };
  })();

  beforeAll(() => {
    Object.defineProperty(window, "localStorage", {
      value: localStorageMock,
    });

    Object.defineProperty(window, "matchMedia", {
      value: jest.fn().mockImplementation((query) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: jest.fn(),
        removeListener: jest.fn(),
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        dispatchEvent: jest.fn(),
      })),
    });
  });

  beforeEach(() => {
    localStorage.clear();
    document.documentElement.classList.remove("dark");
  });

  it("renders the toggle button", () => {
    render(<ThemeToggle />);
    const button = screen.getByRole("button", { name: /switch to dark mode/i });
    expect(button).toBeInTheDocument();
  });

  it("toggles theme on click", () => {
    render(<ThemeToggle />);
    const button = screen.getByRole("button");

    // Initially light mode
    expect(document.documentElement.classList.contains("dark")).toBe(false);
    expect(localStorage.getItem("theme")).toBeNull();

    // Click to enable dark mode
    fireEvent.click(button);
    expect(document.documentElement.classList.contains("dark")).toBe(true);
    expect(localStorage.getItem("theme")).toBe("dark");

    // Click to switch back to light mode
    fireEvent.click(button);
    expect(document.documentElement.classList.contains("dark")).toBe(false);
    expect(localStorage.getItem("theme")).toBe("light");
  });

  it("applies dark mode if saved in localStorage", () => {
    localStorage.setItem("theme", "dark");
    render(<ThemeToggle />);
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });

  it("button aria-label updates correctly", () => {
    render(<ThemeToggle />);
    const button = screen.getByRole("button");
    expect(button).toHaveAttribute("aria-label", "Switch to dark mode");

    fireEvent.click(button);
    expect(button).toHaveAttribute("aria-label", "Switch to light mode");
  });
});
