import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom";
import { MemoryRouter } from "react-router-dom";

console.log(React)

beforeAll(() => {
  Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: jest.fn().mockImplementation(query => ({
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

jest.mock("chart.js", () => ({
  __esModule: true,
  Chart: jest.fn().mockImplementation(() => ({
    destroy: jest.fn(),
  })),
  register: jest.fn(),
  registerables: [],
}));

jest.mock("../src/components/MapModal", () => ({
  __esModule: true,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  default: ({ isOpen }: any) => (
    <div data-testid="map-modal">
      {isOpen && <span>Map Modal Open</span>}
    </div>
  ),
}));

const mockNavigate = jest.fn();
jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useNavigate: () => mockNavigate,
}));

if (!global.fetch) {
  global.fetch = jest.fn();
}

import Dashboard from "../src/pages/Dashboard";

describe("Dashboard Component", () => {
  const mockIntersections = [
    {
      id: "1",
      name: "Intersection 1",
      status: "INTERSECTION_STATUS_OPTIMISED",
      run_count: 2,
      traffic_density: "TRAFFIC_DENSITY_HIGH",
      created_at: "2025-09-26T10:00:00Z",
      details: { address: "123 Main St", city: "CityA", province: "ProvinceA" },
    },
    {
      id: "2",
      name: "Intersection 2",
      status: "unoptimised",
      run_count: 1,
      traffic_density: "TRAFFIC_DENSITY_LOW",
      created_at: "2025-09-27T08:00:00Z",
      details: { address: "456 Side St", city: "CityB", province: "ProvinceB" },
    },
  ];

  beforeEach(() => {
    (global.fetch as jest.Mock).mockImplementation(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ intersections: mockIntersections }),
      })
    );
    localStorage.setItem("authToken", "fake-token");
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test("renders total intersections, active simulations, and optimization runs", async () => {
    render(
      <MemoryRouter>
        <Dashboard />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText("Total Intersections")).toBeInTheDocument();
      const twos = screen.getAllByText("2");
      expect(twos[0]).toBeInTheDocument();
      expect(screen.getByText("Active Simulations")).toBeInTheDocument();
      expect(screen.getByText("3")).toBeInTheDocument();
      expect(screen.getByText("Optimization Runs")).toBeInTheDocument();
      expect(twos[1]).toBeInTheDocument();
    });
  });

  test("renders recent intersections list and allows viewing details", async () => {
  render(
    <MemoryRouter>
      <Dashboard />
    </MemoryRouter>
  );

  await waitFor(() => {
    expect(screen.getByText("Intersection 2")).toBeInTheDocument();
    expect(screen.getByText("Intersection 1")).toBeInTheDocument();
  });

  const viewDetailButtons = screen.getAllByText("View Details");
  fireEvent.click(viewDetailButtons[0]);

  expect(mockNavigate).toHaveBeenCalledWith("/simulation-results", {
    state: expect.objectContaining({
      intersections: ["Intersection 2"],
      intersectionIds: ["2"],
      description: expect.any(String),
      name: expect.any(String),
      type: "simulations",
    }),
  });
});

  test("opens new intersection, run simulation, and map modals", async () => {
    render(
      <MemoryRouter>
        <Dashboard />
      </MemoryRouter>
    );

    const newIntersectionBtn = screen.getByTitle(/Go to Intersections page/i);
    fireEvent.click(newIntersectionBtn);
    expect(mockNavigate).toHaveBeenCalledWith("/intersections");

    const runSimBtn = screen.getByTitle(/Go to Simulations page/i);
    fireEvent.click(runSimBtn);
    expect(mockNavigate).toHaveBeenCalledWith("/simulations");

    const viewMapBtn = screen.getByText(/View Map/i);
    fireEvent.click(viewMapBtn);

    await waitFor(() => {
      expect(screen.getByText(/Map Modal Open/i)).toBeInTheDocument();
    });
  });

  test("handles fetch failure gracefully", async () => {
    (global.fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({ ok: false })
    );

    render(
      <MemoryRouter>
        <Dashboard />
      </MemoryRouter>
    );

    await waitFor(() => {
      const zeros = screen.getAllByText("0");
      expect(zeros[0]).toBeInTheDocument();
    });
  });
});
