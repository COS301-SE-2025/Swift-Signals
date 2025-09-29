import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom";
import { MemoryRouter, Routes, Route } from "react-router-dom";

import SimulationResults, {
  computeStats,
} from "../src/pages/SimulationResults";

console.log(React);

jest.mock("../src/components/Navbar", () => () => <div>Mock Navbar</div>);
jest.mock("../src/components/Footer", () => () => <div>Mock Footer</div>);
jest.mock("../src/components/HelpMenu", () => () => <div>Mock HelpMenu</div>);

describe("computeStats", () => {
  const vehicles = [
    {
      id: "v1",
      positions: [
        { time: 0, x: 0, y: 0, speed: 10 },
        { time: 1, x: 3, y: 4, speed: 20 },
      ],
    },
    {
      id: "v2",
      positions: [
        { time: 0, x: 0, y: 0, speed: 5 },
        { time: 1, x: 0, y: 3, speed: 15 },
      ],
    },
  ];

  it("calculates averages and totals correctly", () => {
    const stats = computeStats(vehicles);
    expect(stats.avgSpeed).toBeCloseTo((10 + 20 + 5 + 15) / 4);
    expect(stats.maxSpeed).toBe(20);
    expect(stats.minSpeed).toBe(5);
    expect(stats.totalDistance).toBeCloseTo(8);
    expect(stats.vehicleCount).toBe(2);
    expect(stats.finalSpeeds).toEqual([20, 15]);
  });
});

describe("SimulationResults Component", () => {
  beforeEach(() => {
    jest.resetAllMocks();
    Storage.prototype.getItem = jest.fn(() => "fake-token");
  });

  function renderWithRoute(path = "/results/123") {
    return render(
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route
            path="/results/:intersectionId"
            element={<SimulationResults />}
          />
        </Routes>
      </MemoryRouter>,
    );
  }

  it("shows loading state initially", () => {
    global.fetch = jest.fn(() => new Promise(() => {})) as jest.Mock;

    renderWithRoute();

    expect(screen.getByText(/Loading/i)).toBeInTheDocument();
  });

  it("shows error state if fetch fails", async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: false,
        statusText: "Not Found",
      } as Response),
    ) as jest.Mock;

    renderWithRoute();

    await waitFor(() =>
      expect(screen.getByText(/Failed to Load Data/i)).toBeInTheDocument(),
    );
  });

  it("renders data state if fetch succeeds", async () => {
    const fakeResponse = {
      output: { vehicles: [] },
      results: {
        average_speed: 10,
        average_travel_time: 5,
        average_waiting_time: 2,
        total_vehicles: 1,
        total_travel_time: 5,
      },
    };

    global.fetch = jest
      .fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => fakeResponse,
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          name: "Test Intersection",
          traffic_density: "LOW",
        }),
      } as Response) as jest.Mock;

    renderWithRoute();

    await waitFor(() =>
      expect(screen.getByText(/Test Intersection/i)).toBeInTheDocument(),
    );
    expect(screen.getByText(/Avg Speed/i)).toBeInTheDocument();
  });
});
