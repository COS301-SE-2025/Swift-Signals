// tests/ComparisonView.test.tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import ComparisonView from "../src/pages/ComparisonView";

console.log(React)

// Mock child components
jest.mock("../src/pages/TrafficSimulation", () => ({
  __esModule: true,
  default: ({ intersectionId, endpoint }: { intersectionId: string; endpoint: string }) => (
    <div data-testid={`TrafficSimulation-${endpoint}-${intersectionId}`}>
      Simulation {endpoint} {intersectionId}
    </div>
  ),
}));

jest.mock("../src/components/HelpMenu", () => ({
  __esModule: true,
  default: () => <div data-testid="HelpMenu">HelpMenu</div>,
}));

// Silence window.close in tests
beforeAll(() => {
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  window.close = jest.fn();
  window.history.back = jest.fn();
});

const renderWithRouter = (state?: any) => {
  return render(
    <MemoryRouter initialEntries={[{ pathname: "/", state }]}>
      <ComparisonView />
    </MemoryRouter>
  );
};

describe("ComparisonView", () => {
  beforeEach(() => {
    localStorage.clear();
    (window.history.back as jest.Mock).mockClear();
    (window.close as jest.Mock).mockClear();
  });

  it("renders original simulation by default", () => {
    renderWithRouter();
    expect(screen.getByText("Original Simulation")).toBeInTheDocument();
    expect(screen.getByTestId("TrafficSimulation-simulate-1")).toBeInTheDocument();
  });

  it("renders optimized simulation if optimizedData is provided", () => {
    renderWithRouter({
      simulationData: {},
      optimizedData: {},
      optimizedIntersectionId: "99",
    });
    expect(screen.getByText("Optimized Simulation")).toBeInTheDocument();
    expect(screen.getByTestId("TrafficSimulation-optimise-1")).toBeInTheDocument();
  });

  it("shows 'No Optimization Available' if optimizedData is missing", () => {
    renderWithRouter({ simulationData: {} });
    expect(screen.getByText("No Optimization Available")).toBeInTheDocument();
    expect(screen.getByText("Check for Optimization")).toBeInTheDocument();
  });

  it("applies dark mode when theme=dark in localStorage", () => {
    localStorage.setItem("theme", "dark");
    renderWithRouter();
    const bodyStyle = window.getComputedStyle(document.body);
    expect(bodyStyle.backgroundColor).toBe("rgb(30, 30, 30)");
  });


  it("exit button triggers history.back when history length > 1", () => {
    // Mock length via a getter
    jest.spyOn(window.history, "length", "get").mockReturnValue(2);
    renderWithRouter();
    fireEvent.click(screen.getByText("Exit"));
    expect(window.history.back).toHaveBeenCalled();
    });

    it("exit button triggers window.close when history length <= 1", () => {
    jest.spyOn(window.history, "length", "get").mockReturnValue(1);
    renderWithRouter();
    fireEvent.click(screen.getByText("Exit"));
    expect(window.close).toHaveBeenCalled();
    });


  it("Escape key triggers exit handler", () => {
    Object.defineProperty(window.history, "length", { value: 2 });
    renderWithRouter();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(window.history.back).toHaveBeenCalled();
  });

  it("toggle left fullscreen button updates text", () => {
    renderWithRouter();
    const button = screen.getByText("Fullscreen");
    fireEvent.click(button);
    expect(screen.getByText("Exit Fullscreen")).toBeInTheDocument();
  });

  it("clicking refresh optimization shows error", () => {
    renderWithRouter({ simulationData: {} });
    fireEvent.click(screen.getByText("Check for Optimization"));
    expect(screen.getByText(/Error: Refresh not available/)).toBeInTheDocument();
  });

  it("renders HelpMenu component", () => {
    renderWithRouter();
    expect(screen.getByTestId("HelpMenu")).toBeInTheDocument();
  });
});
